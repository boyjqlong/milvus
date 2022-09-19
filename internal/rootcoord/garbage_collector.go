package rootcoord

import (
	"context"

	"github.com/milvus-io/milvus/internal/util/typeutil"

	"github.com/milvus-io/milvus/internal/metastore/model"
)

type GarbageCollector interface {
	ReDropCollection(collMeta *model.Collection, ts Timestamp)
	RemoveCreatingCollection(collMeta *model.Collection)
	ReDropPartition(pChannels []string, partition *model.Partition, ts Timestamp)
	GcCollectionData(ctx context.Context, coll *model.Collection, ts typeutil.Timestamp) (ddlTs Timestamp, err error)
	GcPartitionData(ctx context.Context, pChannels []string, partition *model.Partition, ts typeutil.Timestamp) (ddlTs Timestamp, err error)
}

type GarbageCollectorCtx struct {
	s *Core
}

func newGarbageCollectorCtx(s *Core) *GarbageCollectorCtx {
	return &GarbageCollectorCtx{s: s}
}

func (c *GarbageCollectorCtx) ReDropCollection(collMeta *model.Collection, ts Timestamp) {
	// TODO: remove this after data gc can be notified by rpc.
	c.s.chanTimeTick.addDmlChannels(collMeta.PhysicalChannelNames...)

	redo := newBaseRedoTask(c.s.stepExecutor)
	redo.AddAsyncStep(&releaseCollectionStep{
		baseStep:     baseStep{core: c.s},
		collectionID: collMeta.CollectionID,
	})
	redo.AddAsyncStep(&dropIndexStep{
		baseStep: baseStep{core: c.s},
		collID:   collMeta.CollectionID,
	})
	redo.AddAsyncStep(&deleteCollectionDataStep{
		baseStep: baseStep{core: c.s},
		coll:     collMeta,
		ts:       ts,
	})
	redo.AddAsyncStep(&removeDmlChannelsStep{
		baseStep:  baseStep{core: c.s},
		pChannels: collMeta.PhysicalChannelNames,
	})
	redo.AddAsyncStep(&deleteCollectionMetaStep{
		baseStep:     baseStep{core: c.s},
		collectionID: collMeta.CollectionID,
		ts:           ts,
	})

	// err is ignored since no sync steps will be executed.
	_ = redo.Execute(context.Background())
}

func (c *GarbageCollectorCtx) RemoveCreatingCollection(collMeta *model.Collection) {
	redo := newBaseRedoTask(c.s.stepExecutor)
	redo.AddAsyncStep(&unwatchChannelsStep{
		baseStep:     baseStep{core: c.s},
		collectionID: collMeta.CollectionID,
		channels: collectionChannels{
			virtualChannels:  collMeta.VirtualChannelNames,
			physicalChannels: collMeta.PhysicalChannelNames,
		},
	})
	redo.AddAsyncStep(&deleteCollectionMetaStep{
		baseStep:     baseStep{},
		collectionID: collMeta.CollectionID,
		ts:           collMeta.CreateTime,
	})
	// err is ignored since no sync steps will be executed.
	_ = redo.Execute(context.Background())
}

func (c *GarbageCollectorCtx) ReDropPartition(pChannels []string, partition *model.Partition, ts Timestamp) {
	// TODO: remove this after data gc can be notified by rpc.
	c.s.chanTimeTick.addDmlChannels(pChannels...)

	redo := newBaseRedoTask(c.s.stepExecutor)
	redo.AddAsyncStep(&deletePartitionDataStep{
		baseStep:  baseStep{core: c.s},
		pchans:    pChannels,
		partition: partition,
		ts:        ts,
	})
	redo.AddAsyncStep(&removeDmlChannelsStep{
		baseStep:  baseStep{core: c.s},
		pChannels: pChannels,
	})
	redo.AddAsyncStep(&removePartitionMetaStep{
		baseStep:     baseStep{},
		collectionID: partition.CollectionID,
		partitionID:  partition.PartitionID,
		ts:           ts,
	})

	// err is ignored since no sync steps will be executed.
	_ = redo.Execute(context.Background())
}

func (c *GarbageCollectorCtx) GcCollectionData(ctx context.Context, coll *model.Collection, ts typeutil.Timestamp) (ddlTs Timestamp, err error) {
	return c.s.ddlTsLockManager.NotifyCollectionGc(ctx, coll)
}

func (c *GarbageCollectorCtx) GcPartitionData(ctx context.Context, pChannels []string, partition *model.Partition, ts typeutil.Timestamp) (ddlTs Timestamp, err error) {
	return c.s.ddlTsLockManager.NotifyPartitionGc(ctx, pChannels, partition)
}
