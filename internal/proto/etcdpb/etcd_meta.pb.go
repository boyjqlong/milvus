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
	DbName                     string                    `protobuf:"bytes,15,opt,name=db_name,json=dbName,proto3" json:"db_name,omitempty"`
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

func (m *CollectionInfo) GetDbName() string {
	if m != nil {
		return m.DbName
	}
	return ""
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
	// 1037 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xcb, 0x8e, 0xdc, 0x44,
	0x17, 0x8e, 0xdb, 0xd3, 0x17, 0x9f, 0xbe, 0xd7, 0x9f, 0xcc, 0xef, 0x0c, 0x09, 0x38, 0x0d, 0x01,
	0x2b, 0x52, 0x66, 0xc4, 0x0c, 0xb7, 0x0d, 0x88, 0x30, 0x56, 0xa4, 0x16, 0x10, 0xb5, 0x3c, 0xa3,
	0x2c, 0xd8, 0x58, 0xd5, 0x76, 0x4d, 0x77, 0x21, 0xdf, 0xe4, 0xaa, 0x1e, 0x98, 0x37, 0xe0, 0x01,
	0x78, 0x07, 0x1e, 0x81, 0x27, 0xe0, 0x69, 0x58, 0xb3, 0x47, 0x55, 0xe5, 0x7b, 0xf7, 0x20, 0x56,
	0xec, 0x7c, 0xbe, 0xaa, 0x73, 0xea, 0x7c, 0xe7, 0xf2, 0x19, 0xa6, 0x84, 0xfb, 0x81, 0x17, 0x11,
	0x8e, 0x4f, 0xd3, 0x2c, 0xe1, 0x09, 0x9a, 0x47, 0x34, 0xbc, 0xdd, 0x31, 0x65, 0x9d, 0x8a, 0xd3,
	0x93, 0x91, 0x9f, 0x44, 0x51, 0x12, 0x2b, 0xe8, 0x64, 0xc4, 0xfc, 0x2d, 0x89, 0xf2, 0xeb, 0x8b,
	0x3f, 0x34, 0x30, 0x96, 0x71, 0x40, 0x7e, 0x5e, 0xc6, 0x37, 0x09, 0x7a, 0x0a, 0x40, 0x85, 0xe1,
	0xc5, 0x38, 0x22, 0xa6, 0x66, 0x69, 0xb6, 0xe1, 0x1a, 0x12, 0x79, 0x83, 0x23, 0x82, 0x4c, 0xe8,
	0x4b, 0x63, 0xe9, 0x98, 0x1d, 0x4b, 0xb3, 0x75, 0xb7, 0x30, 0x91, 0x03, 0x23, 0xe5, 0x98, 0xe2,
	0x0c, 0x47, 0xcc, 0xd4, 0x2d, 0xdd, 0x1e, 0x9e, 0x3f, 0x3b, 0x6d, 0x24, 0x93, 0xa7, 0xf1, 0x2d,
	0xb9, 0x7b, 0x8b, 0xc3, 0x1d, 0x59, 0x61, 0x9a, 0xb9, 0x43, 0xe9, 0xb6, 0x92, 0x5e, 0x22, 0x7e,
	0x40, 0x42, 0xc2, 0x49, 0x60, 0x1e, 0x59, 0x9a, 0x3d, 0x70, 0x0b, 0x13, 0xbd, 0x07, 0x43, 0x3f,
	0x23, 0x98, 0x13, 0x8f, 0xd3, 0x88, 0x98, 0x5d, 0x4b, 0xb3, 0x8f, 0x5c, 0x50, 0xd0, 0x35, 0x8d,
	0xc8, 0xc2, 0x81, 0xc9, 0x6b, 0x4a, 0xc2, 0xa0, 0xe2, 0x62, 0x42, 0xff, 0x86, 0x86, 0x24, 0x58,
	0x3a, 0x92, 0x88, 0xee, 0x16, 0xe6, 0xfd, 0x34, 0x16, 0xbf, 0xf6, 0x60, 0x72, 0x99, 0x84, 0x21,
	0xf1, 0x39, 0x4d, 0x62, 0x19, 0x66, 0x02, 0x9d, 0x32, 0x42, 0x67, 0xe9, 0xa0, 0x2f, 0xa1, 0xa7,
	0x0a, 0x28, 0x7d, 0x87, 0xe7, 0xcf, 0x9b, 0x1c, 0xf3, 0xe2, 0x56, 0x41, 0xae, 0x24, 0xe0, 0xe6,
	0x4e, 0x6d, 0x22, 0x7a, 0x9b, 0x08, 0x5a, 0xc0, 0x28, 0xc5, 0x19, 0xa7, 0x32, 0x01, 0x87, 0x99,
	0x47, 0x96, 0x6e, 0xeb, 0x6e, 0x03, 0x43, 0x1f, 0xc2, 0xa4, 0xb4, 0x45, 0x63, 0x98, 0xd9, 0xb5,
	0x74, 0xdb, 0x70, 0x5b, 0x28, 0x7a, 0x0d, 0xe3, 0x1b, 0x51, 0x14, 0x4f, 0xf2, 0x23, 0xcc, 0xec,
	0x1d, 0x6a, 0x8b, 0x98, 0x91, 0xd3, 0x66, 0xf1, 0xdc, 0xd1, 0x4d, 0x69, 0x13, 0x86, 0xce, 0xe1,
	0xd1, 0x2d, 0xcd, 0xf8, 0x0e, 0x87, 0x9e, 0xbf, 0xc5, 0x71, 0x4c, 0x42, 0x39, 0x20, 0xcc, 0xec,
	0xcb, 0x67, 0xff, 0x97, 0x1f, 0x5e, 0xaa, 0x33, 0xf5, 0xf6, 0x27, 0x70, 0x9c, 0x6e, 0xef, 0x18,
	0xf5, 0xf7, 0x9c, 0x06, 0xd2, 0xe9, 0x61, 0x71, 0xda, 0xf0, 0xfa, 0x1a, 0x9e, 0x94, 0x1c, 0x3c,
	0x55, 0x95, 0x40, 0x56, 0x8a, 0x71, 0x1c, 0xa5, 0xcc, 0x34, 0x2c, 0xdd, 0x3e, 0x72, 0x4f, 0xca,
	0x3b, 0x97, 0xea, 0xca, 0x75, 0x79, 0x43, 0x8c, 0x30, 0xdb, 0xe2, 0x2c, 0x60, 0x5e, 0xbc, 0x8b,
	0x4c, 0xb0, 0x34, 0xbb, 0xeb, 0x1a, 0x0a, 0x79, 0xb3, 0x8b, 0xd0, 0x12, 0xa6, 0x8c, 0xe3, 0x8c,
	0x7b, 0x69, 0xc2, 0x64, 0x04, 0x66, 0x0e, 0x65, 0x51, 0xac, 0xfb, 0x66, 0xd5, 0xc1, 0x1c, 0xcb,
	0x51, 0x9d, 0x48, 0xc7, 0x55, 0xe1, 0x87, 0x5c, 0x98, 0xfb, 0x49, 0xcc, 0x28, 0xe3, 0x24, 0xf6,
	0xef, 0xbc, 0x90, 0xdc, 0x92, 0xd0, 0x1c, 0x59, 0x9a, 0x3d, 0x69, 0x0f, 0x45, 0x1e, 0xec, 0xb2,
	0xba, 0xfd, 0x9d, 0xb8, 0xec, 0xce, 0xfc, 0x16, 0x82, 0xbe, 0x80, 0x2e, 0xe3, 0x98, 0x13, 0x73,
	0x2c, 0xe3, 0x2c, 0x0e, 0x74, 0xaa, 0x36, 0x5a, 0xe2, 0xa6, 0xab, 0x1c, 0xd0, 0x2b, 0x80, 0x34,
	0x4b, 0x52, 0x92, 0x71, 0x4a, 0x98, 0x39, 0xf9, 0xb7, 0xfb, 0x57, 0x73, 0x42, 0xff, 0x87, 0x7e,
	0xb0, 0x56, 0xab, 0x3f, 0x95, 0xab, 0xdf, 0x0b, 0xd6, 0xa2, 0x2d, 0x8b, 0xbf, 0x34, 0x18, 0xaf,
	0xca, 0x01, 0x14, 0x5b, 0x61, 0xc1, 0xb0, 0x36, 0x91, 0xf9, 0x7a, 0xd4, 0x21, 0xf4, 0x01, 0x8c,
	0x1b, 0xd3, 0x28, 0xd7, 0xc5, 0x70, 0x9b, 0x20, 0xfa, 0x0a, 0xde, 0xf9, 0x87, 0x7e, 0xe7, 0xeb,
	0xf1, 0xf8, 0xde, 0x76, 0xa3, 0xf7, 0x61, 0xec, 0x97, 0xf5, 0xf0, 0xa8, 0xd2, 0x0d, 0xdd, 0x1d,
	0x55, 0xe0, 0x32, 0x40, 0x9f, 0x17, 0x45, 0xed, 0xca, 0xa2, 0x1e, 0x1a, 0xff, 0x92, 0x5d, 0xbd,
	0xa6, 0x8b, 0xdf, 0x34, 0x30, 0x5e, 0x85, 0x14, 0xb3, 0x42, 0x1c, 0xb1, 0x30, 0x1a, 0xe2, 0x28,
	0x11, 0x49, 0x65, 0x2f, 0x95, 0xce, 0x81, 0x54, 0x9e, 0xc1, 0xa8, 0xce, 0x32, 0x27, 0x98, 0x4b,
	0x82, 0xe4, 0x85, 0x2e, 0x8a, 0x6c, 0x8f, 0x64, 0xb6, 0x4f, 0x0f, 0x64, 0x2b, 0x73, 0x6a, 0x64,
	0xfa, 0x4b, 0x07, 0x66, 0x57, 0x64, 0x13, 0x91, 0x98, 0x57, 0x0a, 0xb8, 0x80, 0xfa, 0xe3, 0x45,
	0x97, 0x1a, 0x58, 0xbb, 0x91, 0x9d, 0xfd, 0x46, 0x3e, 0x01, 0x83, 0xe5, 0x91, 0x1d, 0x99, 0xaf,
	0xee, 0x56, 0x80, 0x52, 0x59, 0x21, 0x15, 0x4e, 0x5e, 0xfa, 0xc2, 0xac, 0xab, 0x6c, 0xb7, 0xf9,
	0xb3, 0x30, 0xa1, 0xbf, 0xde, 0x51, 0xe9, 0xd3, 0x53, 0x27, 0xb9, 0x29, 0xca, 0x43, 0x62, 0xbc,
	0x0e, 0x89, 0x52, 0x2c, 0xb3, 0x2f, 0xff, 0x02, 0x43, 0x85, 0x49, 0x62, 0x6d, 0x01, 0x1d, 0xec,
	0xfd, 0x09, 0xfe, 0xd4, 0xea, 0x1a, 0xfe, 0x3d, 0xe1, 0xf8, 0x3f, 0xd7, 0xf0, 0x77, 0x01, 0xca,
	0x0a, 0x15, 0x0a, 0x5e, 0x43, 0xd0, 0xf3, 0x9a, 0x7e, 0x7b, 0x1c, 0x6f, 0x0a, 0xfd, 0xae, 0x96,
	0xe3, 0x1a, 0x6f, 0xd8, 0xde, 0xaf, 0xa0, 0xb7, 0xff, 0x2b, 0x58, 0xfc, 0x2e, 0xd8, 0x66, 0x24,
	0x20, 0x31, 0xa7, 0x38, 0x94, 0x6d, 0x3f, 0x81, 0xc1, 0x8e, 0x91, 0xac, 0x36, 0xa5, 0xa5, 0x8d,
	0x5e, 0x02, 0x22, 0xb1, 0x9f, 0xdd, 0xa5, 0x62, 0x02, 0x53, 0xcc, 0xd8, 0x4f, 0x49, 0x16, 0xe4,
	0xab, 0x39, 0x2f, 0x4f, 0x56, 0xf9, 0x01, 0x3a, 0x86, 0x1e, 0x27, 0x31, 0x8e, 0xb9, 0x24, 0x69,
	0xb8, 0xb9, 0x85, 0x1e, 0xc3, 0x80, 0x32, 0x8f, 0xed, 0x52, 0x92, 0x15, 0x7f, 0x6a, 0xca, 0xae,
	0x84, 0x89, 0x3e, 0x82, 0x29, 0xdb, 0xe2, 0xf3, 0x4f, 0x3f, 0xab, 0xc2, 0x77, 0xa5, 0xef, 0x44,
	0xc1, 0x45, 0xec, 0x17, 0x09, 0x4c, 0x5b, 0x52, 0x86, 0x1e, 0xc1, 0xbc, 0x82, 0xf2, 0x5d, 0x9f,
	0x3d, 0x40, 0xc7, 0x80, 0x5a, 0x30, 0x8d, 0x37, 0x33, 0xad, 0x89, 0x3b, 0x59, 0x92, 0xa6, 0x02,
	0xef, 0x34, 0xc3, 0x48, 0x9c, 0x04, 0x33, 0xfd, 0xc5, 0x8f, 0x30, 0x69, 0xae, 0x39, 0x7a, 0x08,
	0xb3, 0x55, 0x4b, 0x5a, 0x66, 0x0f, 0x84, 0x7b, 0x13, 0x55, 0xaf, 0xd5, 0xe1, 0xda, 0x63, 0xf5,
	0x18, 0xd5, 0x5b, 0x6f, 0x01, 0xaa, 0x25, 0x45, 0x33, 0x18, 0x49, 0xab, 0x7a, 0x63, 0x0e, 0xe3,
	0x0a, 0x51, 0xf1, 0x0b, 0xa8, 0x16, 0xbb, 0xf0, 0x2b, 0xe3, 0x7e, 0x73, 0xf1, 0xc3, 0xc7, 0x1b,
	0xca, 0xb7, 0xbb, 0xb5, 0x10, 0xf3, 0x33, 0x35, 0xb5, 0x2f, 0x69, 0x92, 0x7f, 0x9d, 0xd1, 0x98,
	0x8b, 0x46, 0x87, 0x67, 0x72, 0x90, 0xcf, 0x84, 0x58, 0xa4, 0xeb, 0x75, 0x4f, 0x5a, 0x17, 0x7f,
	0x07, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x86, 0x92, 0xec, 0x2c, 0x0a, 0x00, 0x00,
}
