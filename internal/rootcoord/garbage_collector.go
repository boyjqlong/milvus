package rootcoord

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	pb "github.com/milvus-io/milvus/internal/proto/etcdpb"

	"github.com/milvus-io/milvus/internal/log"

	"github.com/milvus-io/milvus-proto/go-api/commonpb"
	ms "github.com/milvus-io/milvus/internal/mq/msgstream"
	"github.com/milvus-io/milvus/internal/proto/internalpb"
	"github.com/milvus-io/milvus/internal/util/commonpbutil"

	"github.com/milvus-io/milvus/internal/metastore/model"
)

//go:generate mockery --name=GarbageCollector --outpkg=mockrootcoord
type GarbageCollector interface {
	ReDropCollection(collMeta *model.Collection, ts Timestamp)
	RemoveCreatingCollection(collMeta *model.Collection)
	ReDropPartition(pChannels []string, partition *model.Partition, ts Timestamp)
	GcCollectionData(ctx context.Context, coll *model.Collection) (ddlTs Timestamp, err error)
	GcPartitionData(ctx context.Context, pChannels []string, partition *model.Partition) (ddlTs Timestamp, err error)
}

const (
	defaultGCInterval      = time.Second * 3600
	defaultReleaseInterval = time.Second * 60
)

type bgGarbageCollector struct {
	s *Core

	gcInterval      time.Duration
	releaseInterval time.Duration

	shouldBeReleased map[UniqueID]*model.Collection

	wg      sync.WaitGroup
	closeCh chan struct{}

	// TODO(longjiquan): consider partition.
	forceGCChan chan *model.Collection
}

type gcOption func(*bgGarbageCollector)

func (c *bgGarbageCollector) apply(opts ...gcOption) {
	for _, opt := range opts {
		opt(c)
	}
}

func newBgGarbageCollector(s *Core, opts ...gcOption) *bgGarbageCollector {
	c := &bgGarbageCollector{
		s:                s,
		gcInterval:       defaultGCInterval,
		releaseInterval:  defaultReleaseInterval,
		shouldBeReleased: make(map[UniqueID]*model.Collection),
		closeCh:          make(chan struct{}, 1),
		forceGCChan:      make(chan *model.Collection, 1),
	}

	c.apply(opts...)

	return c
}

func (c *bgGarbageCollector) Start() {
	c.wg.Add(1)
	defer c.wg.Done()

	gcTicker := time.NewTicker(c.gcInterval)
	defer gcTicker.Stop()

	releaseTicker := time.NewTicker(c.releaseInterval)
	defer releaseTicker.Stop()

	for {
		select {
		case <-c.closeCh:
			log.Info("garbage collection quit!")
			return
		case coll := <-c.forceGCChan:
			c.handle(coll)
		case <-gcTicker.C:
			c.scan()
		}
	}
}

func (c *bgGarbageCollector) ForceTrigger(coll *model.Collection) {
	c.forceGCChan <- coll
}

func (c *bgGarbageCollector) handle(coll *model.Collection) {
	switch coll.State {
	case pb.CollectionState_CollectionCreated: // available.
		return
	case pb.CollectionState_CollectionCreating: // undo will solve this case.
		return
	case pb.CollectionState_CollectionDropping:
		c.tryReleaseCollection(coll)
	case pb.CollectionState_CollectionDropped: // impossible now.
		log.Warn("garbage collector met a collection whose state is dropped", zap.Int64("collection", coll.CollectionID))
		return

	case pb.CollectionState_CollectionDroppingAndReleased:
		c.tryDropIndex(coll)
	case pb.CollectionState_CollectionDroppingAndIndexDropped:
		c.tryStartGC(coll)
	case pb.CollectionState_CollectionDroppingAndGcStarted:
		c.tryConfirmTsSynced(coll)

	case pb.CollectionState_CollectionDroppingAndTsSynced:
		c.tryConfirmGC(coll)
	case pb.CollectionState_CollectionDroppingAndGcConfirmed:
		c.tryRemoveCollection(coll)
	}
}

type afterAlterCollection func()

func (c *bgGarbageCollector) tryAlterCollection(ctx context.Context, oldColl *model.Collection, newColl *model.Collection, callbacks ...afterAlterCollection) {
	ts, err := c.s.tsoAllocator.GenerateTSO(1)
	if err != nil {
		log.Info("failed to alter collection", zap.Error(err), zap.Int64("collection", oldColl.CollectionID))
		return
	}

	if err := c.s.meta.AlterCollection(ctx, oldColl, newColl, ts); err != nil {
		log.Info("failed to alter collection", zap.Error(err), zap.Int64("collection", oldColl.CollectionID))
		return
	}

	for _, cb := range callbacks {
		cb()
	}
}

func (c *bgGarbageCollector) tryReleaseCollection(coll *model.Collection) {
	ctx := context.Background()

	if err := c.s.broker.ReleaseCollection(ctx, coll.CollectionID); err != nil {
		log.Info("failed to release collection", zap.Error(err), zap.Int64("collection", coll.CollectionID))
		return
	}

	newColl := coll.Clone()
	newColl.State = pb.CollectionState_CollectionDroppingAndReleased

	c.tryAlterCollection(ctx, coll, newColl)
}

func (c *bgGarbageCollector) tryDropIndex(coll *model.Collection) {
	ctx := context.Background()

	if err := c.s.broker.DropCollectionIndex(ctx, coll.CollectionID, nil); err != nil {
		log.Info("failed to drop collection index", zap.Error(err), zap.Int64("collection", coll.CollectionID))
		return
	}

	newColl := coll.Clone()
	newColl.State = pb.CollectionState_CollectionDroppingAndIndexDropped

	c.tryAlterCollection(ctx, coll, newColl)
}

func (c *bgGarbageCollector) tryStartGC(coll *model.Collection) {
	ctx := context.Background()

	ddlTs, err := c.GcCollectionData(ctx, coll)
	if err != nil {
		log.Info("failed to gc collection data", zap.Error(err), zap.Int64("collection", coll.CollectionID))
		return
	}

	newColl := coll.Clone()
	newColl.State = pb.CollectionState_CollectionDroppingAndGcStarted
	newColl.DropTs = ddlTs

	c.tryAlterCollection(ctx, coll, newColl)
}

func (c *bgGarbageCollector) tryConfirmTsSynced(coll *model.Collection) {
	ctx := context.Background()

	for _, channel := range coll.PhysicalChannelNames {
		syncedTs := c.s.chanTimeTick.getSyncedTimeTick(channel)
		if syncedTs < coll.DropTs {
			return
		}
	}

	newColl := coll.Clone()
	newColl.State = pb.CollectionState_CollectionDroppingAndTsSynced

	c.tryAlterCollection(ctx, coll, newColl, func() {
		// memory change should be after persistent change.
		c.s.chanTimeTick.removeDmlChannels(coll.PhysicalChannelNames...)
	})
}

func (c *bgGarbageCollector) tryConfirmGC(coll *model.Collection) {
	ctx := context.Background()

	if !c.s.broker.GcConfirm(ctx, coll.CollectionID, -1) {
		return
	}

	newColl := coll.Clone()
	newColl.State = pb.CollectionState_CollectionDroppingAndGcConfirmed

	c.tryAlterCollection(ctx, coll, newColl)
}

func (c *bgGarbageCollector) tryRemoveCollection(coll *model.Collection) {
	ctx := context.Background()

	ts, err := c.s.tsoAllocator.GenerateTSO(1)
	if err != nil {
		log.Info("failed to remove collection", zap.Error(err), zap.Int64("collection", coll.CollectionID))
		return
	}

	if err := c.s.meta.RemoveCollection(ctx, coll.CollectionID, ts); err != nil {
		log.Info("failed to remove collection", zap.Error(err), zap.Int64("collection", coll.CollectionID))
		return
	}
}

func (c *bgGarbageCollector) scan() {
}

func (c *bgGarbageCollector) Stop() {
	close(c.closeCh)
}

func (c *bgGarbageCollector) ReDropCollection(collMeta *model.Collection, ts Timestamp) {
	// TODO: remove this after data gc can be notified by rpc.
	c.s.chanTimeTick.addDmlChannels(collMeta.PhysicalChannelNames...)
	aliases := c.s.meta.ListAliasesByID(collMeta.CollectionID)

	redo := newBaseRedoTask(c.s.stepExecutor)
	redo.AddAsyncStep(&expireCacheStep{
		baseStep:        baseStep{core: c.s},
		collectionNames: append(aliases, collMeta.Name),
		collectionID:    collMeta.CollectionID,
		ts:              ts,
		opts:            []expireCacheOpt{expireCacheWithDropFlag()},
	})
	redo.AddAsyncStep(&releaseCollectionStep{
		baseStep:     baseStep{core: c.s},
		collectionID: collMeta.CollectionID,
	})
	redo.AddAsyncStep(&dropIndexStep{
		baseStep: baseStep{core: c.s},
		collID:   collMeta.CollectionID,
		partIDs:  nil,
	})
	redo.AddAsyncStep(&deleteCollectionDataStep{
		baseStep: baseStep{core: c.s},
		coll:     collMeta,
	})
	redo.AddAsyncStep(&removeDmlChannelsStep{
		baseStep:  baseStep{core: c.s},
		pChannels: collMeta.PhysicalChannelNames,
	})
	redo.AddAsyncStep(&deleteCollectionMetaStep{
		baseStep:     baseStep{core: c.s},
		collectionID: collMeta.CollectionID,
		// This ts is less than the ts when we notify data nodes to drop collection, but it's OK since we have already
		// marked this collection as deleted. If we want to make this ts greater than the notification's ts, we should
		// wrap a step who will have these three children and connect them with ts.
		ts: ts,
	})

	// err is ignored since no sync steps will be executed.
	_ = redo.Execute(context.Background())
}

func (c *bgGarbageCollector) RemoveCreatingCollection(collMeta *model.Collection) {
	// TODO: remove this after data gc can be notified by rpc.
	c.s.chanTimeTick.addDmlChannels(collMeta.PhysicalChannelNames...)

	redo := newBaseRedoTask(c.s.stepExecutor)

	redo.AddAsyncStep(&unwatchChannelsStep{
		baseStep:     baseStep{core: c.s},
		collectionID: collMeta.CollectionID,
		channels: collectionChannels{
			virtualChannels:  collMeta.VirtualChannelNames,
			physicalChannels: collMeta.PhysicalChannelNames,
		},
	})
	redo.AddAsyncStep(&removeDmlChannelsStep{
		baseStep:  baseStep{core: c.s},
		pChannels: collMeta.PhysicalChannelNames,
	})
	redo.AddAsyncStep(&deleteCollectionMetaStep{
		baseStep:     baseStep{core: c.s},
		collectionID: collMeta.CollectionID,
		// When we undo createCollectionTask, this ts may be less than the ts when unwatch channels.
		ts: collMeta.CreateTime,
	})
	// err is ignored since no sync steps will be executed.
	_ = redo.Execute(context.Background())
}

func (c *bgGarbageCollector) ReDropPartition(pChannels []string, partition *model.Partition, ts Timestamp) {
	// TODO: remove this after data gc can be notified by rpc.
	c.s.chanTimeTick.addDmlChannels(pChannels...)

	redo := newBaseRedoTask(c.s.stepExecutor)
	redo.AddAsyncStep(&expireCacheStep{
		baseStep:     baseStep{core: c.s},
		collectionID: partition.CollectionID,
		ts:           ts,
	})
	redo.AddAsyncStep(&deletePartitionDataStep{
		baseStep:  baseStep{core: c.s},
		pchans:    pChannels,
		partition: partition,
	})
	redo.AddAsyncStep(&removeDmlChannelsStep{
		baseStep:  baseStep{core: c.s},
		pChannels: pChannels,
	})
	redo.AddAsyncStep(&removePartitionMetaStep{
		baseStep:     baseStep{core: c.s},
		collectionID: partition.CollectionID,
		partitionID:  partition.PartitionID,
		// This ts is less than the ts when we notify data nodes to drop partition, but it's OK since we have already
		// marked this partition as deleted. If we want to make this ts greater than the notification's ts, we should
		// wrap a step who will have these children and connect them with ts.
		ts: ts,
	})

	// err is ignored since no sync steps will be executed.
	_ = redo.Execute(context.Background())
}

func (c *bgGarbageCollector) notifyCollectionGc(ctx context.Context, coll *model.Collection) (ddlTs Timestamp, err error) {
	ts, err := c.s.tsoAllocator.GenerateTSO(1)
	if err != nil {
		return 0, err
	}

	msgPack := ms.MsgPack{}
	baseMsg := ms.BaseMsg{
		Ctx:            ctx,
		BeginTimestamp: ts,
		EndTimestamp:   ts,
		HashValues:     []uint32{0},
	}
	msg := &ms.DropCollectionMsg{
		BaseMsg: baseMsg,
		DropCollectionRequest: internalpb.DropCollectionRequest{
			Base: commonpbutil.NewMsgBase(
				commonpbutil.WithMsgType(commonpb.MsgType_DropCollection),
				commonpbutil.WithTimeStamp(ts),
				commonpbutil.WithSourceID(c.s.session.ServerID),
			),
			CollectionName: coll.Name,
			CollectionID:   coll.CollectionID,
		},
	}
	msgPack.Msgs = append(msgPack.Msgs, msg)
	if err := c.s.chanTimeTick.broadcastDmlChannels(coll.PhysicalChannelNames, &msgPack); err != nil {
		return 0, err
	}

	return ts, nil
}

func (c *bgGarbageCollector) notifyPartitionGc(ctx context.Context, pChannels []string, partition *model.Partition) (ddlTs Timestamp, err error) {
	ts, err := c.s.tsoAllocator.GenerateTSO(1)
	if err != nil {
		return 0, err
	}

	msgPack := ms.MsgPack{}
	baseMsg := ms.BaseMsg{
		Ctx:            ctx,
		BeginTimestamp: ts,
		EndTimestamp:   ts,
		HashValues:     []uint32{0},
	}
	msg := &ms.DropPartitionMsg{
		BaseMsg: baseMsg,
		DropPartitionRequest: internalpb.DropPartitionRequest{
			Base: commonpbutil.NewMsgBase(
				commonpbutil.WithMsgType(commonpb.MsgType_DropPartition),
				commonpbutil.WithTimeStamp(ts),
				commonpbutil.WithSourceID(c.s.session.ServerID),
			),
			PartitionName: partition.PartitionName,
			CollectionID:  partition.CollectionID,
			PartitionID:   partition.PartitionID,
		},
	}
	msgPack.Msgs = append(msgPack.Msgs, msg)
	if err := c.s.chanTimeTick.broadcastDmlChannels(pChannels, &msgPack); err != nil {
		return 0, err
	}

	return ts, nil
}

func (c *bgGarbageCollector) GcCollectionData(ctx context.Context, coll *model.Collection) (ddlTs Timestamp, err error) {
	c.s.ddlTsLockManager.Lock()
	c.s.ddlTsLockManager.AddRefCnt(1)
	defer c.s.ddlTsLockManager.AddRefCnt(-1)
	defer c.s.ddlTsLockManager.Unlock()

	ddlTs, err = c.notifyCollectionGc(ctx, coll)
	if err != nil {
		return 0, err
	}
	c.s.ddlTsLockManager.UpdateLastTs(ddlTs)
	return ddlTs, nil
}

func (c *bgGarbageCollector) GcPartitionData(ctx context.Context, pChannels []string, partition *model.Partition) (ddlTs Timestamp, err error) {
	c.s.ddlTsLockManager.Lock()
	c.s.ddlTsLockManager.AddRefCnt(1)
	defer c.s.ddlTsLockManager.AddRefCnt(-1)
	defer c.s.ddlTsLockManager.Unlock()

	ddlTs, err = c.notifyPartitionGc(ctx, pChannels, partition)
	if err != nil {
		return 0, err
	}
	c.s.ddlTsLockManager.UpdateLastTs(ddlTs)
	return ddlTs, nil
}
