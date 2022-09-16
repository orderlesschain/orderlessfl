// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.5
// source: shared.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TargetSystem int32

const (
	TargetSystem_NOTYPESYSTEM      TargetSystem = 0
	TargetSystem_ORDERLESSFL       TargetSystem = 1
	TargetSystem_FEDERATEDLEARNING TargetSystem = 4
)

// Enum value maps for TargetSystem.
var (
	TargetSystem_name = map[int32]string{
		0: "NOTYPESYSTEM",
		1: "ORDERLESSFL",
		4: "FEDERATEDLEARNING",
	}
	TargetSystem_value = map[string]int32{
		"NOTYPESYSTEM":      0,
		"ORDERLESSFL":       1,
		"FEDERATEDLEARNING": 4,
	}
)

func (x TargetSystem) Enum() *TargetSystem {
	p := new(TargetSystem)
	*p = x
	return p
}

func (x TargetSystem) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TargetSystem) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_proto_enumTypes[0].Descriptor()
}

func (TargetSystem) Type() protoreflect.EnumType {
	return &file_shared_proto_enumTypes[0]
}

func (x TargetSystem) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TargetSystem.Descriptor instead.
func (TargetSystem) EnumDescriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{0}
}

type FailureType int32

const (
	FailureType_NOTYPEFAILURE FailureType = 0
	FailureType_CRASHED       FailureType = 1
	FailureType_TAMPERED      FailureType = 2
	FailureType_NOTRESPONDING FailureType = 3
	FailureType_RANODM        FailureType = 4
)

// Enum value maps for FailureType.
var (
	FailureType_name = map[int32]string{
		0: "NOTYPEFAILURE",
		1: "CRASHED",
		2: "TAMPERED",
		3: "NOTRESPONDING",
		4: "RANODM",
	}
	FailureType_value = map[string]int32{
		"NOTYPEFAILURE": 0,
		"CRASHED":       1,
		"TAMPERED":      2,
		"NOTRESPONDING": 3,
		"RANODM":        4,
	}
)

func (x FailureType) Enum() *FailureType {
	p := new(FailureType)
	*p = x
	return p
}

func (x FailureType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FailureType) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_proto_enumTypes[1].Descriptor()
}

func (FailureType) Type() protoreflect.EnumType {
	return &file_shared_proto_enumTypes[1]
}

func (x FailureType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FailureType.Descriptor instead.
func (FailureType) EnumDescriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{1}
}

type Profiling_ProfilingType int32

const (
	Profiling_NONE   Profiling_ProfilingType = 0
	Profiling_CPU    Profiling_ProfilingType = 1
	Profiling_MEMORY Profiling_ProfilingType = 2
)

// Enum value maps for Profiling_ProfilingType.
var (
	Profiling_ProfilingType_name = map[int32]string{
		0: "NONE",
		1: "CPU",
		2: "MEMORY",
	}
	Profiling_ProfilingType_value = map[string]int32{
		"NONE":   0,
		"CPU":    1,
		"MEMORY": 2,
	}
)

func (x Profiling_ProfilingType) Enum() *Profiling_ProfilingType {
	p := new(Profiling_ProfilingType)
	*p = x
	return p
}

func (x Profiling_ProfilingType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Profiling_ProfilingType) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_proto_enumTypes[2].Descriptor()
}

func (Profiling_ProfilingType) Type() protoreflect.EnumType {
	return &file_shared_proto_enumTypes[2]
}

func (x Profiling_ProfilingType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Profiling_ProfilingType.Descriptor instead.
func (Profiling_ProfilingType) EnumDescriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{3, 0}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_shared_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{0}
}

type OperationMode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetSystem                       TargetSystem `protobuf:"varint,1,opt,name=target_system,json=targetSystem,proto3,enum=protos.TargetSystem" json:"target_system,omitempty"`
	EndorsementPolicy                  int32        `protobuf:"varint,12,opt,name=endorsement_policy,json=endorsementPolicy,proto3" json:"endorsement_policy,omitempty"`
	GossipNodeCount                    int32        `protobuf:"varint,3,opt,name=gossip_node_count,json=gossipNodeCount,proto3" json:"gossip_node_count,omitempty"`
	TotalNodeCount                     int32        `protobuf:"varint,4,opt,name=total_node_count,json=totalNodeCount,proto3" json:"total_node_count,omitempty"`
	TotalClientCount                   int32        `protobuf:"varint,5,opt,name=total_client_count,json=totalClientCount,proto3" json:"total_client_count,omitempty"`
	GossipIntervalMs                   int32        `protobuf:"varint,7,opt,name=gossip_interval_ms,json=gossipIntervalMs,proto3" json:"gossip_interval_ms,omitempty"`
	TransactionTimeoutSecond           int32        `protobuf:"varint,8,opt,name=transaction_timeout_second,json=transactionTimeoutSecond,proto3" json:"transaction_timeout_second,omitempty"`
	BlockTimeOutMs                     int32        `protobuf:"varint,10,opt,name=block_time_out_ms,json=blockTimeOutMs,proto3" json:"block_time_out_ms,omitempty"`
	BlockTransactionSize               int32        `protobuf:"varint,11,opt,name=block_transaction_size,json=blockTransactionSize,proto3" json:"block_transaction_size,omitempty"`
	ProposalQueueConsumptionRateTps    int32        `protobuf:"varint,14,opt,name=proposal_queue_consumption_rate_tps,json=proposalQueueConsumptionRateTps,proto3" json:"proposal_queue_consumption_rate_tps,omitempty"`
	TransactionQueueConsumptionRateTps int32        `protobuf:"varint,15,opt,name=transaction_queue_consumption_rate_tps,json=transactionQueueConsumptionRateTps,proto3" json:"transaction_queue_consumption_rate_tps,omitempty"`
	Benchmark                          string       `protobuf:"bytes,16,opt,name=benchmark,proto3" json:"benchmark,omitempty"`
	QueueTickerDurationMs              int32        `protobuf:"varint,17,opt,name=queue_ticker_duration_ms,json=queueTickerDurationMs,proto3" json:"queue_ticker_duration_ms,omitempty"`
	ExtraEndorsementOrgs               int32        `protobuf:"varint,18,opt,name=extra_endorsement_orgs,json=extraEndorsementOrgs,proto3" json:"extra_endorsement_orgs,omitempty"`
	ProfilingEnabled                   string       `protobuf:"bytes,19,opt,name=profiling_enabled,json=profilingEnabled,proto3" json:"profiling_enabled,omitempty"`
	OrgsPercentageIncreasedLoad        int32        `protobuf:"varint,20,opt,name=orgs_percentage_increased_load,json=orgsPercentageIncreasedLoad,proto3" json:"orgs_percentage_increased_load,omitempty"`
	LoadIncreasePercentage             int32        `protobuf:"varint,21,opt,name=load_increase_percentage,json=loadIncreasePercentage,proto3" json:"load_increase_percentage,omitempty"`
}

func (x *OperationMode) Reset() {
	*x = OperationMode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OperationMode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperationMode) ProtoMessage() {}

func (x *OperationMode) ProtoReflect() protoreflect.Message {
	mi := &file_shared_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperationMode.ProtoReflect.Descriptor instead.
func (*OperationMode) Descriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{1}
}

func (x *OperationMode) GetTargetSystem() TargetSystem {
	if x != nil {
		return x.TargetSystem
	}
	return TargetSystem_NOTYPESYSTEM
}

func (x *OperationMode) GetEndorsementPolicy() int32 {
	if x != nil {
		return x.EndorsementPolicy
	}
	return 0
}

func (x *OperationMode) GetGossipNodeCount() int32 {
	if x != nil {
		return x.GossipNodeCount
	}
	return 0
}

func (x *OperationMode) GetTotalNodeCount() int32 {
	if x != nil {
		return x.TotalNodeCount
	}
	return 0
}

func (x *OperationMode) GetTotalClientCount() int32 {
	if x != nil {
		return x.TotalClientCount
	}
	return 0
}

func (x *OperationMode) GetGossipIntervalMs() int32 {
	if x != nil {
		return x.GossipIntervalMs
	}
	return 0
}

func (x *OperationMode) GetTransactionTimeoutSecond() int32 {
	if x != nil {
		return x.TransactionTimeoutSecond
	}
	return 0
}

func (x *OperationMode) GetBlockTimeOutMs() int32 {
	if x != nil {
		return x.BlockTimeOutMs
	}
	return 0
}

func (x *OperationMode) GetBlockTransactionSize() int32 {
	if x != nil {
		return x.BlockTransactionSize
	}
	return 0
}

func (x *OperationMode) GetProposalQueueConsumptionRateTps() int32 {
	if x != nil {
		return x.ProposalQueueConsumptionRateTps
	}
	return 0
}

func (x *OperationMode) GetTransactionQueueConsumptionRateTps() int32 {
	if x != nil {
		return x.TransactionQueueConsumptionRateTps
	}
	return 0
}

func (x *OperationMode) GetBenchmark() string {
	if x != nil {
		return x.Benchmark
	}
	return ""
}

func (x *OperationMode) GetQueueTickerDurationMs() int32 {
	if x != nil {
		return x.QueueTickerDurationMs
	}
	return 0
}

func (x *OperationMode) GetExtraEndorsementOrgs() int32 {
	if x != nil {
		return x.ExtraEndorsementOrgs
	}
	return 0
}

func (x *OperationMode) GetProfilingEnabled() string {
	if x != nil {
		return x.ProfilingEnabled
	}
	return ""
}

func (x *OperationMode) GetOrgsPercentageIncreasedLoad() int32 {
	if x != nil {
		return x.OrgsPercentageIncreasedLoad
	}
	return 0
}

func (x *OperationMode) GetLoadIncreasePercentage() int32 {
	if x != nil {
		return x.LoadIncreasePercentage
	}
	return 0
}

type FailureCommandMode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FailureDurationS int32       `protobuf:"varint,1,opt,name=failure_duration_s,json=failureDurationS,proto3" json:"failure_duration_s,omitempty"`
	FailureType      FailureType `protobuf:"varint,2,opt,name=failure_type,json=failureType,proto3,enum=protos.FailureType" json:"failure_type,omitempty"`
}

func (x *FailureCommandMode) Reset() {
	*x = FailureCommandMode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailureCommandMode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailureCommandMode) ProtoMessage() {}

func (x *FailureCommandMode) ProtoReflect() protoreflect.Message {
	mi := &file_shared_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailureCommandMode.ProtoReflect.Descriptor instead.
func (*FailureCommandMode) Descriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{2}
}

func (x *FailureCommandMode) GetFailureDurationS() int32 {
	if x != nil {
		return x.FailureDurationS
	}
	return 0
}

func (x *FailureCommandMode) GetFailureType() FailureType {
	if x != nil {
		return x.FailureType
	}
	return FailureType_NOTYPEFAILURE
}

type Profiling struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProfilingType Profiling_ProfilingType `protobuf:"varint,1,opt,name=profiling_type,json=profilingType,proto3,enum=protos.Profiling_ProfilingType" json:"profiling_type,omitempty"`
}

func (x *Profiling) Reset() {
	*x = Profiling{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profiling) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profiling) ProtoMessage() {}

func (x *Profiling) ProtoReflect() protoreflect.Message {
	mi := &file_shared_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profiling.ProtoReflect.Descriptor instead.
func (*Profiling) Descriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{3}
}

func (x *Profiling) GetProfilingType() Profiling_ProfilingType {
	if x != nil {
		return x.ProfilingType
	}
	return Profiling_NONE
}

type ProfilingResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ProfilingResult) Reset() {
	*x = ProfilingResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfilingResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfilingResult) ProtoMessage() {}

func (x *ProfilingResult) ProtoReflect() protoreflect.Message {
	mi := &file_shared_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfilingResult.ProtoReflect.Descriptor instead.
func (*ProfilingResult) Descriptor() ([]byte, []int) {
	return file_shared_proto_rawDescGZIP(), []int{4}
}

func (x *ProfilingResult) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_shared_proto protoreflect.FileDescriptor

var file_shared_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0xa5, 0x07, 0x0a, 0x0d, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x6f, 0x64,
	0x65, 0x12, 0x39, 0x0a, 0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x0c,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x2d, 0x0a, 0x12,
	0x65, 0x6e, 0x64, 0x6f, 0x72, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x73,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x2a, 0x0a, 0x11, 0x67,
	0x6f, 0x73, 0x73, 0x69, 0x70, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x4e, 0x6f,
	0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x2c, 0x0a, 0x12, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x2c, 0x0a, 0x12, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x5f, 0x6d, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x67, 0x6f, 0x73,
	0x73, 0x69, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d, 0x73, 0x12, 0x3c, 0x0a,
	0x1a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x18, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x12, 0x29, 0x0a, 0x11, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x6f, 0x75, 0x74, 0x5f, 0x6d, 0x73,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d,
	0x65, 0x4f, 0x75, 0x74, 0x4d, 0x73, 0x12, 0x34, 0x0a, 0x16, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x14, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x4c, 0x0a, 0x23,
	0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x63,
	0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x5f,
	0x74, 0x70, 0x73, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x1f, 0x70, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x61, 0x6c, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x61, 0x74, 0x65, 0x54, 0x70, 0x73, 0x12, 0x52, 0x0a, 0x26, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x61, 0x74, 0x65,
	0x5f, 0x74, 0x70, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x05, 0x52, 0x22, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x74, 0x65, 0x54, 0x70, 0x73, 0x12, 0x1c,
	0x0a, 0x09, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x10, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x37, 0x0a, 0x18,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x5f, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x73, 0x18, 0x11, 0x20, 0x01, 0x28, 0x05, 0x52, 0x15,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x12, 0x34, 0x0a, 0x16, 0x65, 0x78, 0x74, 0x72, 0x61, 0x5f, 0x65,
	0x6e, 0x64, 0x6f, 0x72, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6f, 0x72, 0x67, 0x73, 0x18,
	0x12, 0x20, 0x01, 0x28, 0x05, 0x52, 0x14, 0x65, 0x78, 0x74, 0x72, 0x61, 0x45, 0x6e, 0x64, 0x6f,
	0x72, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4f, 0x72, 0x67, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64,
	0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x69, 0x6e,
	0x67, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x43, 0x0a, 0x1e, 0x6f, 0x72, 0x67, 0x73,
	0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e, 0x63, 0x72,
	0x65, 0x61, 0x73, 0x65, 0x64, 0x5f, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x1b, 0x6f, 0x72, 0x67, 0x73, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65,
	0x49, 0x6e, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x64, 0x4c, 0x6f, 0x61, 0x64, 0x12, 0x38, 0x0a,
	0x18, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x69, 0x6e, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x70,
	0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x16, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6e, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x50, 0x65, 0x72,
	0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x22, 0x7a, 0x0a, 0x12, 0x46, 0x61, 0x69, 0x6c, 0x75,
	0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x2c, 0x0a,
	0x12, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x66, 0x61, 0x69, 0x6c, 0x75,
	0x72, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x12, 0x36, 0x0a, 0x0c, 0x66,
	0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x46, 0x61, 0x69, 0x6c, 0x75,
	0x72, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x22, 0x83, 0x01, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x69, 0x6e,
	0x67, 0x12, 0x46, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2e, 0x0a, 0x0d, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f,
	0x4e, 0x45, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x43, 0x50, 0x55, 0x10, 0x01, 0x12, 0x0a, 0x0a,
	0x06, 0x4d, 0x45, 0x4d, 0x4f, 0x52, 0x59, 0x10, 0x02, 0x22, 0x2b, 0x0a, 0x0f, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2a, 0x48, 0x0a, 0x0c, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x4f, 0x54, 0x59, 0x50, 0x45,
	0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x4f, 0x52, 0x44, 0x45,
	0x52, 0x4c, 0x45, 0x53, 0x53, 0x46, 0x4c, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x46, 0x45, 0x44,
	0x45, 0x52, 0x41, 0x54, 0x45, 0x44, 0x4c, 0x45, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x04,
	0x2a, 0x5a, 0x0a, 0x0b, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x11, 0x0a, 0x0d, 0x4e, 0x4f, 0x54, 0x59, 0x50, 0x45, 0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45,
	0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x52, 0x41, 0x53, 0x48, 0x45, 0x44, 0x10, 0x01, 0x12,
	0x0c, 0x0a, 0x08, 0x54, 0x41, 0x4d, 0x50, 0x45, 0x52, 0x45, 0x44, 0x10, 0x02, 0x12, 0x11, 0x0a,
	0x0d, 0x4e, 0x4f, 0x54, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x03,
	0x12, 0x0a, 0x0a, 0x06, 0x52, 0x41, 0x4e, 0x4f, 0x44, 0x4d, 0x10, 0x04, 0x42, 0x0b, 0x5a, 0x09,
	0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_shared_proto_rawDescOnce sync.Once
	file_shared_proto_rawDescData = file_shared_proto_rawDesc
)

func file_shared_proto_rawDescGZIP() []byte {
	file_shared_proto_rawDescOnce.Do(func() {
		file_shared_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_proto_rawDescData)
	})
	return file_shared_proto_rawDescData
}

var file_shared_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_shared_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_shared_proto_goTypes = []interface{}{
	(TargetSystem)(0),            // 0: protos.TargetSystem
	(FailureType)(0),             // 1: protos.FailureType
	(Profiling_ProfilingType)(0), // 2: protos.Profiling.ProfilingType
	(*Empty)(nil),                // 3: protos.Empty
	(*OperationMode)(nil),        // 4: protos.OperationMode
	(*FailureCommandMode)(nil),   // 5: protos.FailureCommandMode
	(*Profiling)(nil),            // 6: protos.Profiling
	(*ProfilingResult)(nil),      // 7: protos.ProfilingResult
}
var file_shared_proto_depIdxs = []int32{
	0, // 0: protos.OperationMode.target_system:type_name -> protos.TargetSystem
	1, // 1: protos.FailureCommandMode.failure_type:type_name -> protos.FailureType
	2, // 2: protos.Profiling.profiling_type:type_name -> protos.Profiling.ProfilingType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_shared_proto_init() }
func file_shared_proto_init() {
	if File_shared_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_shared_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OperationMode); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_shared_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailureCommandMode); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_shared_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profiling); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_shared_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfilingResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shared_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_proto_goTypes,
		DependencyIndexes: file_shared_proto_depIdxs,
		EnumInfos:         file_shared_proto_enumTypes,
		MessageInfos:      file_shared_proto_msgTypes,
	}.Build()
	File_shared_proto = out.File
	file_shared_proto_rawDesc = nil
	file_shared_proto_goTypes = nil
	file_shared_proto_depIdxs = nil
}
