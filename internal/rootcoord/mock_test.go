package rootcoord

import (
	"context"
	"errors"
	"math/rand"
	"os"

	"github.com/milvus-io/milvus/internal/allocator"
	"github.com/milvus-io/milvus/internal/kv"
	"github.com/milvus-io/milvus/internal/tso"

	"github.com/milvus-io/milvus/internal/log"
	"github.com/milvus-io/milvus/internal/metastore/model"
	"github.com/milvus-io/milvus/internal/proto/commonpb"
	"github.com/milvus-io/milvus/internal/proto/datapb"
	pb "github.com/milvus-io/milvus/internal/proto/etcdpb"
	"github.com/milvus-io/milvus/internal/proto/indexpb"
	"github.com/milvus-io/milvus/internal/proto/internalpb"
	"github.com/milvus-io/milvus/internal/proto/proxypb"
	"github.com/milvus-io/milvus/internal/proto/querypb"
	"github.com/milvus-io/milvus/internal/types"
	"github.com/milvus-io/milvus/internal/util/dependency"
	"github.com/milvus-io/milvus/internal/util/metricsinfo"
	"github.com/milvus-io/milvus/internal/util/retry"
	"github.com/milvus-io/milvus/internal/util/sessionutil"
	"github.com/milvus-io/milvus/internal/util/typeutil"
	"go.uber.org/zap"
)

const (
	TestProxyID     = 100
	TestRootCoordID = 200
)

type mockMetaTable struct {
	IMetaTableV2
	ListCollectionsFunc       func(ctx context.Context, ts Timestamp) ([]*model.Collection, error)
	AddCollectionFunc         func(ctx context.Context, coll *model.Collection) error
	GetCollectionByNameFunc   func(ctx context.Context, collectionName string, ts Timestamp) (*model.Collection, error)
	GetCollectionByIDFunc     func(ctx context.Context, collectionID UniqueID, ts Timestamp) (*model.Collection, error)
	ChangeCollectionStateFunc func(ctx context.Context, collectionID UniqueID, state pb.CollectionState, ts Timestamp) error
	RemoveCollectionFunc      func(ctx context.Context, collectionID UniqueID, ts Timestamp) error
	AddPartitionFunc          func(ctx context.Context, partition *model.Partition) error
	ChangePartitionStateFunc  func(ctx context.Context, collectionID UniqueID, partitionID UniqueID, state pb.PartitionState, ts Timestamp) error
	RemovePartitionFunc       func(ctx context.Context, collectionID UniqueID, partitionID UniqueID, ts Timestamp) error
	CreateAliasFunc           func(ctx context.Context, alias string, collectionName string, ts Timestamp) error
	AlterAliasFunc            func(ctx context.Context, alias string, collectionName string, ts Timestamp) error
	DropAliasFunc             func(ctx context.Context, alias string, ts Timestamp) error
	IsAliasFunc               func(name string) bool
}

func (m *mockMetaTable) ListCollections(ctx context.Context, ts Timestamp) ([]*model.Collection, error) {
	return m.ListCollectionsFunc(ctx, ts)
}

func (m mockMetaTable) AddCollection(ctx context.Context, coll *model.Collection) error {
	return m.AddCollectionFunc(ctx, coll)
}

func (m mockMetaTable) GetCollectionByName(ctx context.Context, collectionName string, ts Timestamp) (*model.Collection, error) {
	return m.GetCollectionByNameFunc(ctx, collectionName, ts)
}

func (m mockMetaTable) GetCollectionByID(ctx context.Context, collectionID UniqueID, ts Timestamp) (*model.Collection, error) {
	return m.GetCollectionByIDFunc(ctx, collectionID, ts)
}

func (m mockMetaTable) ChangeCollectionState(ctx context.Context, collectionID UniqueID, state pb.CollectionState, ts Timestamp) error {
	return m.ChangeCollectionStateFunc(ctx, collectionID, state, ts)
}

func (m mockMetaTable) RemoveCollection(ctx context.Context, collectionID UniqueID, ts Timestamp) error {
	return m.RemoveCollectionFunc(ctx, collectionID, ts)
}

func (m mockMetaTable) AddPartition(ctx context.Context, partition *model.Partition) error {
	return m.AddPartitionFunc(ctx, partition)
}

func (m mockMetaTable) ChangePartitionState(ctx context.Context, collectionID UniqueID, partitionID UniqueID, state pb.PartitionState, ts Timestamp) error {
	return m.ChangePartitionStateFunc(ctx, collectionID, partitionID, state, ts)
}

func (m mockMetaTable) RemovePartition(ctx context.Context, collectionID UniqueID, partitionID UniqueID, ts Timestamp) error {
	return m.RemovePartitionFunc(ctx, collectionID, partitionID, ts)
}

func (m mockMetaTable) CreateAlias(ctx context.Context, alias string, collectionName string, ts Timestamp) error {
	return m.CreateAliasFunc(ctx, alias, collectionName, ts)
}

func (m mockMetaTable) AlterAlias(ctx context.Context, alias string, collectionName string, ts Timestamp) error {
	return m.AlterAliasFunc(ctx, alias, collectionName, ts)
}

func (m mockMetaTable) DropAlias(ctx context.Context, alias string, ts Timestamp) error {
	return m.DropAliasFunc(ctx, alias, ts)
}

func (m mockMetaTable) IsAlias(name string) bool {
	return m.IsAliasFunc(name)
}

func newMockMetaTable() *mockMetaTable {
	return &mockMetaTable{}
}

type mockIndexCoord struct {
	types.IndexCoord
	GetComponentStatesFunc func(ctx context.Context) (*internalpb.ComponentStates, error)
	GetIndexStatesFunc     func(ctx context.Context, req *indexpb.GetIndexStatesRequest) (*indexpb.GetIndexStatesResponse, error)
}

func newMockIndexCoord() *mockIndexCoord {
	return &mockIndexCoord{}
}

func (m *mockIndexCoord) GetComponentStates(ctx context.Context) (*internalpb.ComponentStates, error) {
	return m.GetComponentStatesFunc(ctx)
}

func (m *mockIndexCoord) GetIndexStates(ctx context.Context, req *indexpb.GetIndexStatesRequest) (*indexpb.GetIndexStatesResponse, error) {
	return m.GetIndexStatesFunc(ctx, req)
}

type mockDataCoord struct {
	types.DataCoord
	GetComponentStatesFunc func(ctx context.Context) (*internalpb.ComponentStates, error)
	WatchChannelsFunc      func(ctx context.Context, req *datapb.WatchChannelsRequest) (*datapb.WatchChannelsResponse, error)
	AcquireSegmentLockFunc func(ctx context.Context, req *datapb.AcquireSegmentLockRequest) (*commonpb.Status, error)
	ReleaseSegmentLockFunc func(ctx context.Context, req *datapb.ReleaseSegmentLockRequest) (*commonpb.Status, error)
	FlushFunc              func(ctx context.Context, req *datapb.FlushRequest) (*datapb.FlushResponse, error)
	ImportFunc             func(ctx context.Context, req *datapb.ImportTaskRequest) (*datapb.ImportTaskResponse, error)
}

func newMockDataCoord() *mockDataCoord {
	return &mockDataCoord{}
}

func (m *mockDataCoord) GetComponentStates(ctx context.Context) (*internalpb.ComponentStates, error) {
	return m.GetComponentStatesFunc(ctx)
}

func (m *mockDataCoord) WatchChannels(ctx context.Context, req *datapb.WatchChannelsRequest) (*datapb.WatchChannelsResponse, error) {
	return m.WatchChannelsFunc(ctx, req)
}

func (m *mockDataCoord) AcquireSegmentLock(ctx context.Context, req *datapb.AcquireSegmentLockRequest) (*commonpb.Status, error) {
	return m.AcquireSegmentLockFunc(ctx, req)
}

func (m *mockDataCoord) ReleaseSegmentLock(ctx context.Context, req *datapb.ReleaseSegmentLockRequest) (*commonpb.Status, error) {
	return m.ReleaseSegmentLockFunc(ctx, req)
}

func (m *mockDataCoord) Flush(ctx context.Context, req *datapb.FlushRequest) (*datapb.FlushResponse, error) {
	return m.FlushFunc(ctx, req)
}

func (m *mockDataCoord) Import(ctx context.Context, req *datapb.ImportTaskRequest) (*datapb.ImportTaskResponse, error) {
	return m.ImportFunc(ctx, req)
}

type mockQueryCoord struct {
	types.QueryCoord
	GetSegmentInfoFunc     func(ctx context.Context, req *querypb.GetSegmentInfoRequest) (*querypb.GetSegmentInfoResponse, error)
	GetComponentStatesFunc func(ctx context.Context) (*internalpb.ComponentStates, error)
	ReleaseCollectionFunc  func(ctx context.Context, req *querypb.ReleaseCollectionRequest) (*commonpb.Status, error)
}

func (m mockQueryCoord) GetSegmentInfo(ctx context.Context, req *querypb.GetSegmentInfoRequest) (*querypb.GetSegmentInfoResponse, error) {
	return m.GetSegmentInfoFunc(ctx, req)
}

func (m mockQueryCoord) GetComponentStates(ctx context.Context) (*internalpb.ComponentStates, error) {
	return m.GetComponentStatesFunc(ctx)
}

func (m mockQueryCoord) ReleaseCollection(ctx context.Context, req *querypb.ReleaseCollectionRequest) (*commonpb.Status, error) {
	return m.ReleaseCollectionFunc(ctx, req)
}

func newMockQueryCoord() *mockQueryCoord {
	return &mockQueryCoord{}
}

func newIDAllocator() *allocator.MockGIDAllocator {
	r := allocator.NewMockGIDAllocator()
	r.AllocF = func(count uint32) (allocator.UniqueID, allocator.UniqueID, error) {
		return 0, 0, nil
	}
	r.AllocOneF = func() (allocator.UniqueID, error) {
		return 0, nil
	}
	return r
}

func newTsoAllocator() *tso.MockAllocator {
	r := tso.NewMockAllocator()
	r.GenerateTSOF = func(count uint32) (uint64, error) {
		return 0, nil
	}
	return r
}

func newTxnKV() *kv.TxnKVMock {
	r := kv.NewMockTxnKV()
	r.SaveF = func(key, value string) error {
		return nil
	}
	r.RemoveF = func(key string) error {
		return nil
	}
	return r
}

type mockProxy struct {
	types.Proxy
	InvalidateCollectionMetaCacheFunc func(ctx context.Context, request *proxypb.InvalidateCollMetaCacheRequest) (*commonpb.Status, error)
}

func (m mockProxy) InvalidateCollectionMetaCache(ctx context.Context, request *proxypb.InvalidateCollMetaCacheRequest) (*commonpb.Status, error) {
	return m.InvalidateCollectionMetaCacheFunc(ctx, request)
}

func newMockProxy() *mockProxy {
	r := &mockProxy{}
	r.InvalidateCollectionMetaCacheFunc = func(ctx context.Context, request *proxypb.InvalidateCollMetaCacheRequest) (*commonpb.Status, error) {
		return succStatus(), nil
	}
	return r
}

func newTestCore(opts ...Opt) *RootCoord {
	c := &RootCoord{
		session: &sessionutil.Session{ServerID: TestRootCoordID},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func withValidProxyManager() Opt {
	return func(c *RootCoord) {
		c.proxyClientManager = &proxyClientManager{
			proxyClient: make(map[UniqueID]types.Proxy),
		}
		p := newMockProxy()
		p.InvalidateCollectionMetaCacheFunc = func(ctx context.Context, request *proxypb.InvalidateCollMetaCacheRequest) (*commonpb.Status, error) {
			return succStatus(), nil
		}
	}
}

func withInvalidProxyManager() Opt {
	return func(c *RootCoord) {
		c.proxyClientManager = &proxyClientManager{
			proxyClient: make(map[UniqueID]types.Proxy),
		}
		p := newMockProxy()
		p.InvalidateCollectionMetaCacheFunc = func(ctx context.Context, request *proxypb.InvalidateCollMetaCacheRequest) (*commonpb.Status, error) {
			return succStatus(), errors.New("error mock InvalidateCollectionMetaCache")
		}
		c.proxyClientManager.proxyClient[TestProxyID] = p
	}
}

func withMeta(meta IMetaTableV2) Opt {
	return func(c *RootCoord) {
		c.meta = meta
	}
}

func withInvalidMeta() Opt {
	meta := newMockMetaTable()
	meta.ListCollectionsFunc = func(ctx context.Context, ts Timestamp) ([]*model.Collection, error) {
		return nil, errors.New("error mock ListCollections")
	}
	meta.GetCollectionByNameFunc = func(ctx context.Context, collectionName string, ts Timestamp) (*model.Collection, error) {
		return nil, errors.New("error mock GetCollectionByName")
	}
	meta.GetCollectionByIDFunc = func(ctx context.Context, collectionID typeutil.UniqueID, ts Timestamp) (*model.Collection, error) {
		return nil, errors.New("error mock GetCollectionByID")
	}
	meta.AddPartitionFunc = func(ctx context.Context, partition *model.Partition) error {
		return errors.New("error mock AddPartition")
	}
	meta.ChangePartitionStateFunc = func(ctx context.Context, collectionID UniqueID, partitionID UniqueID, state pb.PartitionState, ts Timestamp) error {
		return errors.New("error mock ChangePartitionState")
	}
	meta.CreateAliasFunc = func(ctx context.Context, alias string, collectionName string, ts Timestamp) error {
		return errors.New("error mock CreateAlias")
	}
	meta.AlterAliasFunc = func(ctx context.Context, alias string, collectionName string, ts Timestamp) error {
		return errors.New("error mock AlterAlias")
	}
	meta.DropAliasFunc = func(ctx context.Context, alias string, ts Timestamp) error {
		return errors.New("error mock DropAlias")
	}
	return withMeta(meta)
}

func withIdAllocator(idAllocator allocator.GIDAllocator) Opt {
	return func(c *RootCoord) {
		c.idAllocator = idAllocator
	}
}

func withValidIdAllocator() Opt {
	idAllocator := newIDAllocator()
	idAllocator.AllocOneF = func() (allocator.UniqueID, error) {
		return rand.Int63(), nil
	}
	return withIdAllocator(idAllocator)
}

func withInvalidIdAllocator() Opt {
	idAllocator := newIDAllocator()
	idAllocator.AllocOneF = func() (allocator.UniqueID, error) {
		return -1, errors.New("error mock AllocOne")
	}
	idAllocator.AllocF = func(count uint32) (allocator.UniqueID, allocator.UniqueID, error) {
		return -1, -1, errors.New("error mock Alloc")
	}
	return withIdAllocator(idAllocator)
}

func withQueryCoord(qc types.QueryCoord) Opt {
	return func(c *RootCoord) {
		c.queryCoord = qc
	}
}

func withUnhealthyQueryCoord() Opt {
	qc := newMockQueryCoord()
	qc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Abnormal},
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "error mock GetComponentStates"),
		}, retry.Unrecoverable(errors.New("error mock GetComponentStates"))
	}
	return withQueryCoord(qc)
}

func withInvalidQueryCoord() Opt {
	qc := newMockQueryCoord()
	qc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	qc.ReleaseCollectionFunc = func(ctx context.Context, req *querypb.ReleaseCollectionRequest) (*commonpb.Status, error) {
		return nil, errors.New("error mock ReleaseCollection")
	}
	qc.GetSegmentInfoFunc = func(ctx context.Context, req *querypb.GetSegmentInfoRequest) (*querypb.GetSegmentInfoResponse, error) {
		return nil, errors.New("error mock GetSegmentInfo")
	}
	return withQueryCoord(qc)
}

func withFailedQueryCoord() Opt {
	qc := newMockQueryCoord()
	qc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	qc.ReleaseCollectionFunc = func(ctx context.Context, req *querypb.ReleaseCollectionRequest) (*commonpb.Status, error) {
		return failStatus(commonpb.ErrorCode_UnexpectedError, "mock release collection error"), nil
	}
	qc.GetSegmentInfoFunc = func(ctx context.Context, req *querypb.GetSegmentInfoRequest) (*querypb.GetSegmentInfoResponse, error) {
		return &querypb.GetSegmentInfoResponse{
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "mock get segment info error"),
		}, nil
	}
	return withQueryCoord(qc)
}

func withValidQueryCoord() Opt {
	qc := newMockQueryCoord()
	qc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	qc.ReleaseCollectionFunc = func(ctx context.Context, req *querypb.ReleaseCollectionRequest) (*commonpb.Status, error) {
		return succStatus(), nil
	}
	qc.GetSegmentInfoFunc = func(ctx context.Context, req *querypb.GetSegmentInfoRequest) (*querypb.GetSegmentInfoResponse, error) {
		return &querypb.GetSegmentInfoResponse{
			Status: succStatus(),
		}, nil
	}
	return withQueryCoord(qc)
}

func withIndexCoord(ic types.IndexCoord) Opt {
	return func(c *RootCoord) {
		c.indexCoord = ic
	}
}

func withUnhealthyIndexCoord() Opt {
	ic := newMockIndexCoord()
	ic.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Abnormal},
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "error mock GetComponentStates"),
		}, retry.Unrecoverable(errors.New("error mock GetComponentStates"))
	}
	return withIndexCoord(ic)
}

func withInvalidIndexCoord() Opt {
	ic := newMockIndexCoord()
	ic.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	ic.GetIndexStatesFunc = func(ctx context.Context, req *indexpb.GetIndexStatesRequest) (*indexpb.GetIndexStatesResponse, error) {
		return nil, errors.New("error mock GetIndexStates")
	}
	return withIndexCoord(ic)
}

func withFailedIndexCoord() Opt {
	ic := newMockIndexCoord()
	ic.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	ic.GetIndexStatesFunc = func(ctx context.Context, req *indexpb.GetIndexStatesRequest) (*indexpb.GetIndexStatesResponse, error) {
		return &indexpb.GetIndexStatesResponse{
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "mock get index states error"),
		}, nil
	}
	return withIndexCoord(ic)
}

func withValidIndexCoord() Opt {
	ic := newMockIndexCoord()
	ic.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	ic.GetIndexStatesFunc = func(ctx context.Context, req *indexpb.GetIndexStatesRequest) (*indexpb.GetIndexStatesResponse, error) {
		return &indexpb.GetIndexStatesResponse{
			Status: succStatus(),
		}, nil
	}
	return withIndexCoord(ic)
}

// cleanTestEnv clean test environment, for example, files generated by rocksmq.
func cleanTestEnv() {
	path := "/tmp/milvus"
	if err := os.RemoveAll(path); err != nil {
		log.Warn("failed to clean test directories", zap.Error(err), zap.String("path", path))
	}
	log.Debug("clean test environment", zap.String("path", path))
}

func withTtSynchronizer(ticker *timetickSync) Opt {
	return func(c *RootCoord) {
		c.chanTimeTick = ticker
	}
}

func newRocksMqTtSynchronizer() *timetickSync {
	Params.InitOnce()
	Params.RootCoordCfg.DmlChannelNum = 4
	ctx := context.Background()
	factory := dependency.NewDefaultFactory(true)
	chans := map[UniqueID][]string{}
	ticker := newTimeTickSync(ctx, TestRootCoordID, factory, chans)
	return ticker
}

// cleanTestEnv should be called if tested with this option.
func withRocksMqTtSynchronizer() Opt {
	ticker := newRocksMqTtSynchronizer()
	return withTtSynchronizer(ticker)
}

func withDataCoord(dc types.DataCoord) Opt {
	return func(c *RootCoord) {
		c.dataCoord = dc
	}
}

func withUnhealthyDataCoord() Opt {
	dc := newMockDataCoord()
	dc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Abnormal},
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "error mock GetComponentStates"),
		}, retry.Unrecoverable(errors.New("error mock GetComponentStates"))
	}
	return withDataCoord(dc)
}

func withInvalidDataCoord() Opt {
	dc := newMockDataCoord()
	dc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	dc.WatchChannelsFunc = func(ctx context.Context, req *datapb.WatchChannelsRequest) (*datapb.WatchChannelsResponse, error) {
		return nil, errors.New("error mock WatchChannels")
	}
	dc.WatchChannelsFunc = func(ctx context.Context, req *datapb.WatchChannelsRequest) (*datapb.WatchChannelsResponse, error) {
		return nil, errors.New("error mock WatchChannels")
	}
	dc.AcquireSegmentLockFunc = func(ctx context.Context, req *datapb.AcquireSegmentLockRequest) (*commonpb.Status, error) {
		return nil, errors.New("error mock AddSegRefLock")
	}
	dc.ReleaseSegmentLockFunc = func(ctx context.Context, req *datapb.ReleaseSegmentLockRequest) (*commonpb.Status, error) {
		return nil, errors.New("error mock ReleaseSegRefLock")
	}
	dc.FlushFunc = func(ctx context.Context, req *datapb.FlushRequest) (*datapb.FlushResponse, error) {
		return nil, errors.New("error mock Flush")
	}
	dc.ImportFunc = func(ctx context.Context, req *datapb.ImportTaskRequest) (*datapb.ImportTaskResponse, error) {
		return nil, errors.New("error mock Import")
	}
	return withDataCoord(dc)
}

func withFailedDataCoord() Opt {
	dc := newMockDataCoord()
	dc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	dc.WatchChannelsFunc = func(ctx context.Context, req *datapb.WatchChannelsRequest) (*datapb.WatchChannelsResponse, error) {
		return &datapb.WatchChannelsResponse{
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "mock watch channels error"),
		}, nil
	}
	dc.AcquireSegmentLockFunc = func(ctx context.Context, req *datapb.AcquireSegmentLockRequest) (*commonpb.Status, error) {
		return failStatus(commonpb.ErrorCode_UnexpectedError, "mock add seg ref lock error"), nil
	}
	dc.ReleaseSegmentLockFunc = func(ctx context.Context, req *datapb.ReleaseSegmentLockRequest) (*commonpb.Status, error) {
		return failStatus(commonpb.ErrorCode_UnexpectedError, "mock release seg ref lock error"), nil
	}
	dc.FlushFunc = func(ctx context.Context, req *datapb.FlushRequest) (*datapb.FlushResponse, error) {
		return &datapb.FlushResponse{
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "mock flush error"),
		}, nil
	}
	dc.ImportFunc = func(ctx context.Context, req *datapb.ImportTaskRequest) (*datapb.ImportTaskResponse, error) {
		return &datapb.ImportTaskResponse{
			Status: failStatus(commonpb.ErrorCode_UnexpectedError, "mock import error"),
		}, nil
	}
	return withDataCoord(dc)
}

func withValidDataCoord() Opt {
	dc := newMockDataCoord()
	dc.GetComponentStatesFunc = func(ctx context.Context) (*internalpb.ComponentStates, error) {
		return &internalpb.ComponentStates{
			State:  &internalpb.ComponentInfo{StateCode: internalpb.StateCode_Healthy},
			Status: succStatus(),
		}, nil
	}
	dc.WatchChannelsFunc = func(ctx context.Context, req *datapb.WatchChannelsRequest) (*datapb.WatchChannelsResponse, error) {
		return &datapb.WatchChannelsResponse{
			Status: succStatus(),
		}, nil
	}
	dc.AcquireSegmentLockFunc = func(ctx context.Context, req *datapb.AcquireSegmentLockRequest) (*commonpb.Status, error) {
		return succStatus(), nil
	}
	dc.ReleaseSegmentLockFunc = func(ctx context.Context, req *datapb.ReleaseSegmentLockRequest) (*commonpb.Status, error) {
		return succStatus(), nil
	}
	dc.FlushFunc = func(ctx context.Context, req *datapb.FlushRequest) (*datapb.FlushResponse, error) {
		return &datapb.FlushResponse{
			Status: succStatus(),
		}, nil
	}
	dc.ImportFunc = func(ctx context.Context, req *datapb.ImportTaskRequest) (*datapb.ImportTaskResponse, error) {
		return &datapb.ImportTaskResponse{
			Status: succStatus(),
		}, nil
	}
	return withDataCoord(dc)
}

func withStateCode(code internalpb.StateCode) Opt {
	return func(c *RootCoord) {
		c.UpdateStateCode(code)
	}
}

func withHealthyCode() Opt {
	return withStateCode(internalpb.StateCode_Healthy)
}

func withAbnormalCode() Opt {
	return withStateCode(internalpb.StateCode_Abnormal)
}

type mockScheduler struct {
	IScheduler
	AddTaskFunc func(t taskV2) error
}

func newMockScheduler() *mockScheduler {
	return &mockScheduler{}
}

func (m mockScheduler) AddTask(t taskV2) error {
	if m.AddTaskFunc != nil {
		return m.AddTaskFunc(t)
	}
	return nil
}

func withScheduler(sched IScheduler) Opt {
	return func(c *RootCoord) {
		c.scheduler = sched
	}
}

func withValidScheduler() Opt {
	sched := newMockScheduler()
	sched.AddTaskFunc = func(t taskV2) error {
		t.NotifyDone(nil)
		return nil
	}
	return withScheduler(sched)
}

func withInvalidScheduler() Opt {
	sched := newMockScheduler()
	sched.AddTaskFunc = func(t taskV2) error {
		return errors.New("error mock AddTask")
	}
	return withScheduler(sched)
}

func withTaskFailScheduler() Opt {
	sched := newMockScheduler()
	sched.AddTaskFunc = func(t taskV2) error {
		err := errors.New("error mock task fail")
		t.NotifyDone(err)
		return nil
	}
	return withScheduler(sched)
}

func withTsoAllocator(alloc tso.Allocator) Opt {
	return func(c *RootCoord) {
		c.tsoAllocator = alloc
	}
}

func withInvalidTsoAllocator() Opt {
	alloc := newTsoAllocator()
	alloc.GenerateTSOF = func(count uint32) (uint64, error) {
		return 0, errors.New("error mock GenerateTSO")
	}
	return withTsoAllocator(alloc)
}

func withMetricsCacheManager() Opt {
	return func(c *RootCoord) {
		m := metricsinfo.NewMetricsCacheManager()
		c.metricsCacheManager = m
	}
}

type mockBroker struct {
	Broker

	ReleaseCollectionFunc   func(ctx context.Context, collectionID UniqueID) error
	GetQuerySegmentInfoFunc func(ctx context.Context, collectionID int64, segIDs []int64) (retResp *querypb.GetSegmentInfoResponse, retErr error)

	WatchChannelsFunc     func(ctx context.Context, info *watchInfo) error
	UnwatchChannelsFunc   func(ctx context.Context, info *watchInfo) error
	AddSegRefLockFunc     func(ctx context.Context, taskID int64, segIDs []int64) error
	ReleaseSegRefLockFunc func(ctx context.Context, taskID int64, segIDs []int64) error
	FlushFunc             func(ctx context.Context, cID int64, segIDs []int64) error
	ImportFunc            func(ctx context.Context, req *datapb.ImportTaskRequest) (*datapb.ImportTaskResponse, error)

	GetIndexStatesFunc func(ctx context.Context, IndexBuildIDs []int64) (idxInfo []*indexpb.IndexInfo, retErr error)
}

func newMockBroker() *mockBroker {
	return &mockBroker{}
}

func (b mockBroker) WatchChannels(ctx context.Context, info *watchInfo) error {
	return b.WatchChannelsFunc(ctx, info)
}

func (b mockBroker) UnwatchChannels(ctx context.Context, info *watchInfo) error {
	return b.UnwatchChannelsFunc(ctx, info)
}

func withBroker(b Broker) Opt {
	return func(c *RootCoord) {
		c.broker = b
	}
}
