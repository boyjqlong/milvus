// Code generated by protoc-gen-go. DO NOT EDIT.
// source: etcd_meta.proto

package etcdpb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	commonpb "github.com/milvus-io/milvus-proto/go-api/commonpb"
	schemapb "github.com/milvus-io/milvus-proto/go-api/schemapb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type CollectionState int32

const (
	CollectionState_CollectionCreated  CollectionState = 0
	CollectionState_CollectionCreating CollectionState = 1
	CollectionState_CollectionDropping CollectionState = 2
	CollectionState_CollectionDropped  CollectionState = 3
)

var CollectionState_name = map[int32]string{
	0: "CollectionCreated",
	1: "CollectionCreating",
	2: "CollectionDropping",
	3: "CollectionDropped",
}

var CollectionState_value = map[string]int32{
	"CollectionCreated":  0,
	"CollectionCreating": 1,
	"CollectionDropping": 2,
	"CollectionDropped":  3,
}

func (x CollectionState) String() string {
	return proto.EnumName(CollectionState_name, int32(x))
}

func (CollectionState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{0}
}

type PartitionState int32

const (
	PartitionState_PartitionCreated  PartitionState = 0
	PartitionState_PartitionCreating PartitionState = 1
	PartitionState_PartitionDropping PartitionState = 2
	PartitionState_PartitionDropped  PartitionState = 3
)

var PartitionState_name = map[int32]string{
	0: "PartitionCreated",
	1: "PartitionCreating",
	2: "PartitionDropping",
	3: "PartitionDropped",
}

var PartitionState_value = map[string]int32{
	"PartitionCreated":  0,
	"PartitionCreating": 1,
	"PartitionDropping": 2,
	"PartitionDropped":  3,
}

func (x PartitionState) String() string {
	return proto.EnumName(PartitionState_name, int32(x))
}

func (PartitionState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{1}
}

type AliasState int32

const (
	AliasState_AliasCreated  AliasState = 0
	AliasState_AliasCreating AliasState = 1
	AliasState_AliasDropping AliasState = 2
	AliasState_AliasDropped  AliasState = 3
)

var AliasState_name = map[int32]string{
	0: "AliasCreated",
	1: "AliasCreating",
	2: "AliasDropping",
	3: "AliasDropped",
}

var AliasState_value = map[string]int32{
	"AliasCreated":  0,
	"AliasCreating": 1,
	"AliasDropping": 2,
	"AliasDropped":  3,
}

func (x AliasState) String() string {
	return proto.EnumName(AliasState_name, int32(x))
}

func (AliasState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{2}
}

type IndexInfo struct {
	IndexName            string                   `protobuf:"bytes,1,opt,name=index_name,json=indexName,proto3" json:"index_name,omitempty"`
	IndexID              int64                    `protobuf:"varint,2,opt,name=indexID,proto3" json:"indexID,omitempty"`
	IndexParams          []*commonpb.KeyValuePair `protobuf:"bytes,3,rep,name=index_params,json=indexParams,proto3" json:"index_params,omitempty"`
	Deleted              bool                     `protobuf:"varint,4,opt,name=deleted,proto3" json:"deleted,omitempty"`
	CreateTime           uint64                   `protobuf:"varint,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *IndexInfo) Reset()         { *m = IndexInfo{} }
func (m *IndexInfo) String() string { return proto.CompactTextString(m) }
func (*IndexInfo) ProtoMessage()    {}
func (*IndexInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{0}
}

func (m *IndexInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IndexInfo.Unmarshal(m, b)
}
func (m *IndexInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IndexInfo.Marshal(b, m, deterministic)
}
func (m *IndexInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexInfo.Merge(m, src)
}
func (m *IndexInfo) XXX_Size() int {
	return xxx_messageInfo_IndexInfo.Size(m)
}
func (m *IndexInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexInfo.DiscardUnknown(m)
}

var xxx_messageInfo_IndexInfo proto.InternalMessageInfo

func (m *IndexInfo) GetIndexName() string {
	if m != nil {
		return m.IndexName
	}
	return ""
}

func (m *IndexInfo) GetIndexID() int64 {
	if m != nil {
		return m.IndexID
	}
	return 0
}

func (m *IndexInfo) GetIndexParams() []*commonpb.KeyValuePair {
	if m != nil {
		return m.IndexParams
	}
	return nil
}

func (m *IndexInfo) GetDeleted() bool {
	if m != nil {
		return m.Deleted
	}
	return false
}

func (m *IndexInfo) GetCreateTime() uint64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

type FieldIndexInfo struct {
	FiledID              int64    `protobuf:"varint,1,opt,name=filedID,proto3" json:"filedID,omitempty"`
	IndexID              int64    `protobuf:"varint,2,opt,name=indexID,proto3" json:"indexID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FieldIndexInfo) Reset()         { *m = FieldIndexInfo{} }
func (m *FieldIndexInfo) String() string { return proto.CompactTextString(m) }
func (*FieldIndexInfo) ProtoMessage()    {}
func (*FieldIndexInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{1}
}

func (m *FieldIndexInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldIndexInfo.Unmarshal(m, b)
}
func (m *FieldIndexInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldIndexInfo.Marshal(b, m, deterministic)
}
func (m *FieldIndexInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldIndexInfo.Merge(m, src)
}
func (m *FieldIndexInfo) XXX_Size() int {
	return xxx_messageInfo_FieldIndexInfo.Size(m)
}
func (m *FieldIndexInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldIndexInfo.DiscardUnknown(m)
}

var xxx_messageInfo_FieldIndexInfo proto.InternalMessageInfo

func (m *FieldIndexInfo) GetFiledID() int64 {
	if m != nil {
		return m.FiledID
	}
	return 0
}

func (m *FieldIndexInfo) GetIndexID() int64 {
	if m != nil {
		return m.IndexID
	}
	return 0
}

type CollectionInfo struct {
	ID         int64                      `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Schema     *schemapb.CollectionSchema `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
	CreateTime uint64                     `protobuf:"varint,3,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// deprecate
	PartitionIDs []int64 `protobuf:"varint,4,rep,packed,name=partitionIDs,proto3" json:"partitionIDs,omitempty"`
	// deprecate
	PartitionNames []string `protobuf:"bytes,5,rep,name=partitionNames,proto3" json:"partitionNames,omitempty"`
	// deprecate
	FieldIndexes         []*FieldIndexInfo `protobuf:"bytes,6,rep,name=field_indexes,json=fieldIndexes,proto3" json:"field_indexes,omitempty"`
	VirtualChannelNames  []string          `protobuf:"bytes,7,rep,name=virtual_channel_names,json=virtualChannelNames,proto3" json:"virtual_channel_names,omitempty"`
	PhysicalChannelNames []string          `protobuf:"bytes,8,rep,name=physical_channel_names,json=physicalChannelNames,proto3" json:"physical_channel_names,omitempty"`
	// deprecate
	PartitionCreatedTimestamps []uint64                  `protobuf:"varint,9,rep,packed,name=partition_created_timestamps,json=partitionCreatedTimestamps,proto3" json:"partition_created_timestamps,omitempty"`
	ShardsNum                  int32                     `protobuf:"varint,10,opt,name=shards_num,json=shardsNum,proto3" json:"shards_num,omitempty"`
	StartPositions             []*commonpb.KeyDataPair   `protobuf:"bytes,11,rep,name=start_positions,json=startPositions,proto3" json:"start_positions,omitempty"`
	ConsistencyLevel           commonpb.ConsistencyLevel `protobuf:"varint,12,opt,name=consistency_level,json=consistencyLevel,proto3,enum=milvus.proto.common.ConsistencyLevel" json:"consistency_level,omitempty"`
	State                      CollectionState           `protobuf:"varint,13,opt,name=state,proto3,enum=milvus.proto.etcd.CollectionState" json:"state,omitempty"`
	Properties                 []*commonpb.KeyValuePair  `protobuf:"bytes,14,rep,name=properties,proto3" json:"properties,omitempty"`
	DropTs                     uint64                    `protobuf:"varint,15,opt,name=drop_ts,json=dropTs,proto3" json:"drop_ts,omitempty"`
	XXX_NoUnkeyedLiteral       struct{}                  `json:"-"`
	XXX_unrecognized           []byte                    `json:"-"`
	XXX_sizecache              int32                     `json:"-"`
}

func (m *CollectionInfo) Reset()         { *m = CollectionInfo{} }
func (m *CollectionInfo) String() string { return proto.CompactTextString(m) }
func (*CollectionInfo) ProtoMessage()    {}
func (*CollectionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{2}
}

func (m *CollectionInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollectionInfo.Unmarshal(m, b)
}
func (m *CollectionInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollectionInfo.Marshal(b, m, deterministic)
}
func (m *CollectionInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollectionInfo.Merge(m, src)
}
func (m *CollectionInfo) XXX_Size() int {
	return xxx_messageInfo_CollectionInfo.Size(m)
}
func (m *CollectionInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_CollectionInfo.DiscardUnknown(m)
}

var xxx_messageInfo_CollectionInfo proto.InternalMessageInfo

func (m *CollectionInfo) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *CollectionInfo) GetSchema() *schemapb.CollectionSchema {
	if m != nil {
		return m.Schema
	}
	return nil
}

func (m *CollectionInfo) GetCreateTime() uint64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *CollectionInfo) GetPartitionIDs() []int64 {
	if m != nil {
		return m.PartitionIDs
	}
	return nil
}

func (m *CollectionInfo) GetPartitionNames() []string {
	if m != nil {
		return m.PartitionNames
	}
	return nil
}

func (m *CollectionInfo) GetFieldIndexes() []*FieldIndexInfo {
	if m != nil {
		return m.FieldIndexes
	}
	return nil
}

func (m *CollectionInfo) GetVirtualChannelNames() []string {
	if m != nil {
		return m.VirtualChannelNames
	}
	return nil
}

func (m *CollectionInfo) GetPhysicalChannelNames() []string {
	if m != nil {
		return m.PhysicalChannelNames
	}
	return nil
}

func (m *CollectionInfo) GetPartitionCreatedTimestamps() []uint64 {
	if m != nil {
		return m.PartitionCreatedTimestamps
	}
	return nil
}

func (m *CollectionInfo) GetShardsNum() int32 {
	if m != nil {
		return m.ShardsNum
	}
	return 0
}

func (m *CollectionInfo) GetStartPositions() []*commonpb.KeyDataPair {
	if m != nil {
		return m.StartPositions
	}
	return nil
}

func (m *CollectionInfo) GetConsistencyLevel() commonpb.ConsistencyLevel {
	if m != nil {
		return m.ConsistencyLevel
	}
	return commonpb.ConsistencyLevel_Strong
}

func (m *CollectionInfo) GetState() CollectionState {
	if m != nil {
		return m.State
	}
	return CollectionState_CollectionCreated
}

func (m *CollectionInfo) GetProperties() []*commonpb.KeyValuePair {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *CollectionInfo) GetDropTs() uint64 {
	if m != nil {
		return m.DropTs
	}
	return 0
}

type PartitionInfo struct {
	PartitionID               int64          `protobuf:"varint,1,opt,name=partitionID,proto3" json:"partitionID,omitempty"`
	PartitionName             string         `protobuf:"bytes,2,opt,name=partitionName,proto3" json:"partitionName,omitempty"`
	PartitionCreatedTimestamp uint64         `protobuf:"varint,3,opt,name=partition_created_timestamp,json=partitionCreatedTimestamp,proto3" json:"partition_created_timestamp,omitempty"`
	CollectionId              int64          `protobuf:"varint,4,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"`
	State                     PartitionState `protobuf:"varint,5,opt,name=state,proto3,enum=milvus.proto.etcd.PartitionState" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}       `json:"-"`
	XXX_unrecognized          []byte         `json:"-"`
	XXX_sizecache             int32          `json:"-"`
}

func (m *PartitionInfo) Reset()         { *m = PartitionInfo{} }
func (m *PartitionInfo) String() string { return proto.CompactTextString(m) }
func (*PartitionInfo) ProtoMessage()    {}
func (*PartitionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{3}
}

func (m *PartitionInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PartitionInfo.Unmarshal(m, b)
}
func (m *PartitionInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PartitionInfo.Marshal(b, m, deterministic)
}
func (m *PartitionInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PartitionInfo.Merge(m, src)
}
func (m *PartitionInfo) XXX_Size() int {
	return xxx_messageInfo_PartitionInfo.Size(m)
}
func (m *PartitionInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PartitionInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PartitionInfo proto.InternalMessageInfo

func (m *PartitionInfo) GetPartitionID() int64 {
	if m != nil {
		return m.PartitionID
	}
	return 0
}

func (m *PartitionInfo) GetPartitionName() string {
	if m != nil {
		return m.PartitionName
	}
	return ""
}

func (m *PartitionInfo) GetPartitionCreatedTimestamp() uint64 {
	if m != nil {
		return m.PartitionCreatedTimestamp
	}
	return 0
}

func (m *PartitionInfo) GetCollectionId() int64 {
	if m != nil {
		return m.CollectionId
	}
	return 0
}

func (m *PartitionInfo) GetState() PartitionState {
	if m != nil {
		return m.State
	}
	return PartitionState_PartitionCreated
}

type AliasInfo struct {
	AliasName            string     `protobuf:"bytes,1,opt,name=alias_name,json=aliasName,proto3" json:"alias_name,omitempty"`
	CollectionId         int64      `protobuf:"varint,2,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"`
	CreatedTime          uint64     `protobuf:"varint,3,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	State                AliasState `protobuf:"varint,4,opt,name=state,proto3,enum=milvus.proto.etcd.AliasState" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AliasInfo) Reset()         { *m = AliasInfo{} }
func (m *AliasInfo) String() string { return proto.CompactTextString(m) }
func (*AliasInfo) ProtoMessage()    {}
func (*AliasInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{4}
}

func (m *AliasInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AliasInfo.Unmarshal(m, b)
}
func (m *AliasInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AliasInfo.Marshal(b, m, deterministic)
}
func (m *AliasInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AliasInfo.Merge(m, src)
}
func (m *AliasInfo) XXX_Size() int {
	return xxx_messageInfo_AliasInfo.Size(m)
}
func (m *AliasInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AliasInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AliasInfo proto.InternalMessageInfo

func (m *AliasInfo) GetAliasName() string {
	if m != nil {
		return m.AliasName
	}
	return ""
}

func (m *AliasInfo) GetCollectionId() int64 {
	if m != nil {
		return m.CollectionId
	}
	return 0
}

func (m *AliasInfo) GetCreatedTime() uint64 {
	if m != nil {
		return m.CreatedTime
	}
	return 0
}

func (m *AliasInfo) GetState() AliasState {
	if m != nil {
		return m.State
	}
	return AliasState_AliasCreated
}

type SegmentIndexInfo struct {
	CollectionID         int64    `protobuf:"varint,1,opt,name=collectionID,proto3" json:"collectionID,omitempty"`
	PartitionID          int64    `protobuf:"varint,2,opt,name=partitionID,proto3" json:"partitionID,omitempty"`
	SegmentID            int64    `protobuf:"varint,3,opt,name=segmentID,proto3" json:"segmentID,omitempty"`
	FieldID              int64    `protobuf:"varint,4,opt,name=fieldID,proto3" json:"fieldID,omitempty"`
	IndexID              int64    `protobuf:"varint,5,opt,name=indexID,proto3" json:"indexID,omitempty"`
	BuildID              int64    `protobuf:"varint,6,opt,name=buildID,proto3" json:"buildID,omitempty"`
	EnableIndex          bool     `protobuf:"varint,7,opt,name=enable_index,json=enableIndex,proto3" json:"enable_index,omitempty"`
	CreateTime           uint64   `protobuf:"varint,8,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SegmentIndexInfo) Reset()         { *m = SegmentIndexInfo{} }
func (m *SegmentIndexInfo) String() string { return proto.CompactTextString(m) }
func (*SegmentIndexInfo) ProtoMessage()    {}
func (*SegmentIndexInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{5}
}

func (m *SegmentIndexInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SegmentIndexInfo.Unmarshal(m, b)
}
func (m *SegmentIndexInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SegmentIndexInfo.Marshal(b, m, deterministic)
}
func (m *SegmentIndexInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SegmentIndexInfo.Merge(m, src)
}
func (m *SegmentIndexInfo) XXX_Size() int {
	return xxx_messageInfo_SegmentIndexInfo.Size(m)
}
func (m *SegmentIndexInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SegmentIndexInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SegmentIndexInfo proto.InternalMessageInfo

func (m *SegmentIndexInfo) GetCollectionID() int64 {
	if m != nil {
		return m.CollectionID
	}
	return 0
}

func (m *SegmentIndexInfo) GetPartitionID() int64 {
	if m != nil {
		return m.PartitionID
	}
	return 0
}

func (m *SegmentIndexInfo) GetSegmentID() int64 {
	if m != nil {
		return m.SegmentID
	}
	return 0
}

func (m *SegmentIndexInfo) GetFieldID() int64 {
	if m != nil {
		return m.FieldID
	}
	return 0
}

func (m *SegmentIndexInfo) GetIndexID() int64 {
	if m != nil {
		return m.IndexID
	}
	return 0
}

func (m *SegmentIndexInfo) GetBuildID() int64 {
	if m != nil {
		return m.BuildID
	}
	return 0
}

func (m *SegmentIndexInfo) GetEnableIndex() bool {
	if m != nil {
		return m.EnableIndex
	}
	return false
}

func (m *SegmentIndexInfo) GetCreateTime() uint64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

// TODO move to proto files of interprocess communication
type CollectionMeta struct {
	ID                   int64                      `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Schema               *schemapb.CollectionSchema `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
	CreateTime           uint64                     `protobuf:"varint,3,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	SegmentIDs           []int64                    `protobuf:"varint,4,rep,packed,name=segmentIDs,proto3" json:"segmentIDs,omitempty"`
	PartitionTags        []string                   `protobuf:"bytes,5,rep,name=partition_tags,json=partitionTags,proto3" json:"partition_tags,omitempty"`
	PartitionIDs         []int64                    `protobuf:"varint,6,rep,packed,name=partitionIDs,proto3" json:"partitionIDs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *CollectionMeta) Reset()         { *m = CollectionMeta{} }
func (m *CollectionMeta) String() string { return proto.CompactTextString(m) }
func (*CollectionMeta) ProtoMessage()    {}
func (*CollectionMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{6}
}

func (m *CollectionMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollectionMeta.Unmarshal(m, b)
}
func (m *CollectionMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollectionMeta.Marshal(b, m, deterministic)
}
func (m *CollectionMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollectionMeta.Merge(m, src)
}
func (m *CollectionMeta) XXX_Size() int {
	return xxx_messageInfo_CollectionMeta.Size(m)
}
func (m *CollectionMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_CollectionMeta.DiscardUnknown(m)
}

var xxx_messageInfo_CollectionMeta proto.InternalMessageInfo

func (m *CollectionMeta) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *CollectionMeta) GetSchema() *schemapb.CollectionSchema {
	if m != nil {
		return m.Schema
	}
	return nil
}

func (m *CollectionMeta) GetCreateTime() uint64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *CollectionMeta) GetSegmentIDs() []int64 {
	if m != nil {
		return m.SegmentIDs
	}
	return nil
}

func (m *CollectionMeta) GetPartitionTags() []string {
	if m != nil {
		return m.PartitionTags
	}
	return nil
}

func (m *CollectionMeta) GetPartitionIDs() []int64 {
	if m != nil {
		return m.PartitionIDs
	}
	return nil
}

type CredentialInfo struct {
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// encrypted by bcrypt (for higher security level)
	EncryptedPassword string `protobuf:"bytes,2,opt,name=encrypted_password,json=encryptedPassword,proto3" json:"encrypted_password,omitempty"`
	Tenant            string `protobuf:"bytes,3,opt,name=tenant,proto3" json:"tenant,omitempty"`
	IsSuper           bool   `protobuf:"varint,4,opt,name=is_super,json=isSuper,proto3" json:"is_super,omitempty"`
	// encrypted by sha256 (for good performance in cache mapping)
	Sha256Password       string   `protobuf:"bytes,5,opt,name=sha256_password,json=sha256Password,proto3" json:"sha256_password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CredentialInfo) Reset()         { *m = CredentialInfo{} }
func (m *CredentialInfo) String() string { return proto.CompactTextString(m) }
func (*CredentialInfo) ProtoMessage()    {}
func (*CredentialInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_975d306d62b73e88, []int{7}
}

func (m *CredentialInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CredentialInfo.Unmarshal(m, b)
}
func (m *CredentialInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CredentialInfo.Marshal(b, m, deterministic)
}
func (m *CredentialInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CredentialInfo.Merge(m, src)
}
func (m *CredentialInfo) XXX_Size() int {
	return xxx_messageInfo_CredentialInfo.Size(m)
}
func (m *CredentialInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_CredentialInfo.DiscardUnknown(m)
}

var xxx_messageInfo_CredentialInfo proto.InternalMessageInfo

func (m *CredentialInfo) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *CredentialInfo) GetEncryptedPassword() string {
	if m != nil {
		return m.EncryptedPassword
	}
	return ""
}

func (m *CredentialInfo) GetTenant() string {
	if m != nil {
		return m.Tenant
	}
	return ""
}

func (m *CredentialInfo) GetIsSuper() bool {
	if m != nil {
		return m.IsSuper
	}
	return false
}

func (m *CredentialInfo) GetSha256Password() string {
	if m != nil {
		return m.Sha256Password
	}
	return ""
}

func init() {
	proto.RegisterEnum("milvus.proto.etcd.CollectionState", CollectionState_name, CollectionState_value)
	proto.RegisterEnum("milvus.proto.etcd.PartitionState", PartitionState_name, PartitionState_value)
	proto.RegisterEnum("milvus.proto.etcd.AliasState", AliasState_name, AliasState_value)
	proto.RegisterType((*IndexInfo)(nil), "milvus.proto.etcd.IndexInfo")
	proto.RegisterType((*FieldIndexInfo)(nil), "milvus.proto.etcd.FieldIndexInfo")
	proto.RegisterType((*CollectionInfo)(nil), "milvus.proto.etcd.CollectionInfo")
	proto.RegisterType((*PartitionInfo)(nil), "milvus.proto.etcd.PartitionInfo")
	proto.RegisterType((*AliasInfo)(nil), "milvus.proto.etcd.AliasInfo")
	proto.RegisterType((*SegmentIndexInfo)(nil), "milvus.proto.etcd.SegmentIndexInfo")
	proto.RegisterType((*CollectionMeta)(nil), "milvus.proto.etcd.CollectionMeta")
	proto.RegisterType((*CredentialInfo)(nil), "milvus.proto.etcd.CredentialInfo")
}

func init() { proto.RegisterFile("etcd_meta.proto", fileDescriptor_975d306d62b73e88) }

var fileDescriptor_975d306d62b73e88 = []byte{
	// 1039 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xcb, 0x6e, 0xe4, 0x44,
	0x17, 0x1e, 0xb7, 0xd3, 0xdd, 0xf1, 0xe9, 0x4b, 0x3a, 0xf5, 0xcf, 0xe4, 0xaf, 0x09, 0x33, 0xe0,
	0x69, 0x18, 0xb0, 0x46, 0x9a, 0x44, 0x24, 0xdc, 0x36, 0x20, 0x86, 0x58, 0x23, 0xb5, 0x80, 0x51,
	0xcb, 0x89, 0x66, 0xc1, 0xc6, 0xaa, 0xd8, 0x95, 0xee, 0x42, 0xbe, 0xc9, 0x55, 0x1d, 0xc8, 0x1b,
	0xf0, 0x00, 0xbc, 0x03, 0x8f, 0xc0, 0x13, 0xf0, 0x34, 0xac, 0xd9, 0xa3, 0xaa, 0xf2, 0xbd, 0x3b,
	0x88, 0x15, 0x3b, 0x7f, 0x5f, 0xd5, 0x39, 0x75, 0xae, 0x9f, 0xe1, 0x80, 0x8a, 0x20, 0xf4, 0x63,
	0x2a, 0xc8, 0x49, 0x96, 0xa7, 0x22, 0x45, 0x87, 0x31, 0x8b, 0x6e, 0x37, 0x5c, 0xa3, 0x13, 0x79,
	0x7a, 0x3c, 0x0e, 0xd2, 0x38, 0x4e, 0x13, 0x4d, 0x1d, 0x8f, 0x79, 0xb0, 0xa6, 0x71, 0x71, 0x7d,
	0xfe, 0x87, 0x01, 0xd6, 0x22, 0x09, 0xe9, 0xcf, 0x8b, 0xe4, 0x26, 0x45, 0x4f, 0x01, 0x98, 0x04,
	0x7e, 0x42, 0x62, 0x8a, 0x0d, 0xdb, 0x70, 0x2c, 0xcf, 0x52, 0xcc, 0x1b, 0x12, 0x53, 0x84, 0x61,
	0xa8, 0xc0, 0xc2, 0xc5, 0x3d, 0xdb, 0x70, 0x4c, 0xaf, 0x84, 0xc8, 0x85, 0xb1, 0x36, 0xcc, 0x48,
	0x4e, 0x62, 0x8e, 0x4d, 0xdb, 0x74, 0x46, 0x67, 0xcf, 0x4e, 0x5a, 0xc1, 0x14, 0x61, 0x7c, 0x4b,
	0xef, 0xde, 0x92, 0x68, 0x43, 0x97, 0x84, 0xe5, 0xde, 0x48, 0x99, 0x2d, 0x95, 0x95, 0xf4, 0x1f,
	0xd2, 0x88, 0x0a, 0x1a, 0xe2, 0x3d, 0xdb, 0x70, 0xf6, 0xbd, 0x12, 0xa2, 0xf7, 0x60, 0x14, 0xe4,
	0x94, 0x08, 0xea, 0x0b, 0x16, 0x53, 0xdc, 0xb7, 0x0d, 0x67, 0xcf, 0x03, 0x4d, 0x5d, 0xb1, 0x98,
	0xce, 0x5d, 0x98, 0xbe, 0x66, 0x34, 0x0a, 0xeb, 0x5c, 0x30, 0x0c, 0x6f, 0x58, 0x44, 0xc3, 0x85,
	0xab, 0x12, 0x31, 0xbd, 0x12, 0xde, 0x9f, 0xc6, 0xfc, 0xd7, 0x01, 0x4c, 0x2f, 0xd2, 0x28, 0xa2,
	0x81, 0x60, 0x69, 0xa2, 0xdc, 0x4c, 0xa1, 0x57, 0x79, 0xe8, 0x2d, 0x5c, 0xf4, 0x25, 0x0c, 0x74,
	0x01, 0x95, 0xed, 0xe8, 0xec, 0x79, 0x3b, 0xc7, 0xa2, 0xb8, 0xb5, 0x93, 0x4b, 0x45, 0x78, 0x85,
	0x51, 0x37, 0x11, 0xb3, 0x9b, 0x08, 0x9a, 0xc3, 0x38, 0x23, 0xb9, 0x60, 0x2a, 0x00, 0x97, 0xe3,
	0x3d, 0xdb, 0x74, 0x4c, 0xaf, 0xc5, 0xa1, 0x0f, 0x61, 0x5a, 0x61, 0xd9, 0x18, 0x8e, 0xfb, 0xb6,
	0xe9, 0x58, 0x5e, 0x87, 0x45, 0xaf, 0x61, 0x72, 0x23, 0x8b, 0xe2, 0xab, 0xfc, 0x28, 0xc7, 0x83,
	0x5d, 0x6d, 0x91, 0x33, 0x72, 0xd2, 0x2e, 0x9e, 0x37, 0xbe, 0xa9, 0x30, 0xe5, 0xe8, 0x0c, 0x1e,
	0xdd, 0xb2, 0x5c, 0x6c, 0x48, 0xe4, 0x07, 0x6b, 0x92, 0x24, 0x34, 0x52, 0x03, 0xc2, 0xf1, 0x50,
	0x3d, 0xfb, 0xbf, 0xe2, 0xf0, 0x42, 0x9f, 0xe9, 0xb7, 0x3f, 0x81, 0xa3, 0x6c, 0x7d, 0xc7, 0x59,
	0xb0, 0x65, 0xb4, 0xaf, 0x8c, 0x1e, 0x96, 0xa7, 0x2d, 0xab, 0xaf, 0xe1, 0x49, 0x95, 0x83, 0xaf,
	0xab, 0x12, 0xaa, 0x4a, 0x71, 0x41, 0xe2, 0x8c, 0x63, 0xcb, 0x36, 0x9d, 0x3d, 0xef, 0xb8, 0xba,
	0x73, 0xa1, 0xaf, 0x5c, 0x55, 0x37, 0xe4, 0x08, 0xf3, 0x35, 0xc9, 0x43, 0xee, 0x27, 0x9b, 0x18,
	0x83, 0x6d, 0x38, 0x7d, 0xcf, 0xd2, 0xcc, 0x9b, 0x4d, 0x8c, 0x16, 0x70, 0xc0, 0x05, 0xc9, 0x85,
	0x9f, 0xa5, 0x5c, 0x79, 0xe0, 0x78, 0xa4, 0x8a, 0x62, 0xdf, 0x37, 0xab, 0x2e, 0x11, 0x44, 0x8d,
	0xea, 0x54, 0x19, 0x2e, 0x4b, 0x3b, 0xe4, 0xc1, 0x61, 0x90, 0x26, 0x9c, 0x71, 0x41, 0x93, 0xe0,
	0xce, 0x8f, 0xe8, 0x2d, 0x8d, 0xf0, 0xd8, 0x36, 0x9c, 0x69, 0x77, 0x28, 0x0a, 0x67, 0x17, 0xf5,
	0xed, 0xef, 0xe4, 0x65, 0x6f, 0x16, 0x74, 0x18, 0xf4, 0x05, 0xf4, 0xb9, 0x20, 0x82, 0xe2, 0x89,
	0xf2, 0x33, 0xdf, 0xd1, 0xa9, 0xc6, 0x68, 0xc9, 0x9b, 0x9e, 0x36, 0x40, 0xaf, 0x00, 0xb2, 0x3c,
	0xcd, 0x68, 0x2e, 0x18, 0xe5, 0x78, 0xfa, 0x6f, 0xf7, 0xaf, 0x61, 0x84, 0xfe, 0x0f, 0xc3, 0x30,
	0x4f, 0x33, 0x5f, 0x70, 0x7c, 0xa0, 0xe6, 0x72, 0x20, 0xe1, 0x15, 0x9f, 0xff, 0x65, 0xc0, 0x64,
	0x59, 0x0d, 0xa0, 0xdc, 0x0a, 0x1b, 0x46, 0x8d, 0x89, 0x2c, 0xd6, 0xa3, 0x49, 0xa1, 0x0f, 0x60,
	0xd2, 0x9a, 0x46, 0xb5, 0x2e, 0x96, 0xd7, 0x26, 0xd1, 0x57, 0xf0, 0xce, 0x3f, 0xf4, 0xbb, 0x58,
	0x8f, 0xc7, 0xf7, 0xb6, 0x1b, 0xbd, 0x0f, 0x93, 0xa0, 0xaa, 0x87, 0xcf, 0xb4, 0x6e, 0x98, 0xde,
	0xb8, 0x26, 0x17, 0x21, 0xfa, 0xbc, 0x2c, 0x6a, 0x5f, 0x15, 0x75, 0xd7, 0xf8, 0x57, 0xd9, 0x35,
	0x6b, 0x3a, 0xff, 0xcd, 0x00, 0xeb, 0x55, 0xc4, 0x08, 0x2f, 0xc5, 0x91, 0x48, 0xd0, 0x12, 0x47,
	0xc5, 0xa8, 0x54, 0xb6, 0x42, 0xe9, 0xed, 0x08, 0xe5, 0x19, 0x8c, 0x9b, 0x59, 0x16, 0x09, 0x16,
	0x92, 0xa0, 0xf2, 0x42, 0xe7, 0x65, 0xb4, 0x7b, 0x2a, 0xda, 0xa7, 0x3b, 0xa2, 0x55, 0x31, 0xb5,
	0x22, 0xfd, 0xa5, 0x07, 0xb3, 0x4b, 0xba, 0x8a, 0x69, 0x22, 0x6a, 0x05, 0x9c, 0x43, 0xf3, 0xf1,
	0xb2, 0x4b, 0x2d, 0xae, 0xdb, 0xc8, 0xde, 0x76, 0x23, 0x9f, 0x80, 0xc5, 0x0b, 0xcf, 0xae, 0x8a,
	0xd7, 0xf4, 0x6a, 0x42, 0xab, 0xac, 0x94, 0x0a, 0xb7, 0x28, 0x7d, 0x09, 0x9b, 0x2a, 0xdb, 0x6f,
	0xff, 0x2c, 0x30, 0x0c, 0xaf, 0x37, 0x4c, 0xd9, 0x0c, 0xf4, 0x49, 0x01, 0x65, 0x79, 0x68, 0x42,
	0xae, 0x23, 0xaa, 0x15, 0x0b, 0x0f, 0xd5, 0x5f, 0x60, 0xa4, 0x39, 0x95, 0x58, 0x57, 0x40, 0xf7,
	0xb7, 0xfe, 0x04, 0x7f, 0x1a, 0x4d, 0x0d, 0xff, 0x9e, 0x0a, 0xf2, 0x9f, 0x6b, 0xf8, 0xbb, 0x00,
	0x55, 0x85, 0x4a, 0x05, 0x6f, 0x30, 0xe8, 0x79, 0x43, 0xbf, 0x7d, 0x41, 0x56, 0xa5, 0x7e, 0xd7,
	0xcb, 0x71, 0x45, 0x56, 0x7c, 0xeb, 0x57, 0x30, 0xd8, 0xfe, 0x15, 0xcc, 0x7f, 0x97, 0xd9, 0xe6,
	0x34, 0xa4, 0x89, 0x60, 0x24, 0x52, 0x6d, 0x3f, 0x86, 0xfd, 0x0d, 0xa7, 0x79, 0x63, 0x4a, 0x2b,
	0x8c, 0x5e, 0x02, 0xa2, 0x49, 0x90, 0xdf, 0x65, 0x72, 0x02, 0x33, 0xc2, 0xf9, 0x4f, 0x69, 0x1e,
	0x16, 0xab, 0x79, 0x58, 0x9d, 0x2c, 0x8b, 0x03, 0x74, 0x04, 0x03, 0x41, 0x13, 0x92, 0x08, 0x95,
	0xa4, 0xe5, 0x15, 0x08, 0x3d, 0x86, 0x7d, 0xc6, 0x7d, 0xbe, 0xc9, 0x68, 0x5e, 0xfe, 0xa9, 0x19,
	0xbf, 0x94, 0x10, 0x7d, 0x04, 0x07, 0x7c, 0x4d, 0xce, 0x3e, 0xfd, 0xac, 0x76, 0xdf, 0x57, 0xb6,
	0x53, 0x4d, 0x97, 0xbe, 0x5f, 0xa4, 0x70, 0xd0, 0x91, 0x32, 0xf4, 0x08, 0x0e, 0x6b, 0xaa, 0xd8,
	0xf5, 0xd9, 0x03, 0x74, 0x04, 0xa8, 0x43, 0xb3, 0x64, 0x35, 0x33, 0xda, 0xbc, 0x9b, 0xa7, 0x59,
	0x26, 0xf9, 0x5e, 0xdb, 0x8d, 0xe2, 0x69, 0x38, 0x33, 0x5f, 0xfc, 0x08, 0xd3, 0xf6, 0x9a, 0xa3,
	0x87, 0x30, 0x5b, 0x76, 0xa4, 0x65, 0xf6, 0x40, 0x9a, 0xb7, 0x59, 0xfd, 0x5a, 0x93, 0x6e, 0x3c,
	0xd6, 0xf4, 0x51, 0xbf, 0xf5, 0x16, 0xa0, 0x5e, 0x52, 0x34, 0x83, 0xb1, 0x42, 0xf5, 0x1b, 0x87,
	0x30, 0xa9, 0x19, 0xed, 0xbf, 0xa4, 0x1a, 0xbe, 0x4b, 0xbb, 0xca, 0xef, 0x37, 0xe7, 0x3f, 0x7c,
	0xbc, 0x62, 0x62, 0xbd, 0xb9, 0x96, 0x62, 0x7e, 0xaa, 0xa7, 0xf6, 0x25, 0x4b, 0x8b, 0xaf, 0x53,
	0x96, 0x08, 0xd9, 0xe8, 0xe8, 0x54, 0x0d, 0xf2, 0xa9, 0x14, 0x8b, 0xec, 0xfa, 0x7a, 0xa0, 0xd0,
	0xf9, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x5d, 0x04, 0x5f, 0x2c, 0x0a, 0x00, 0x00,
}
