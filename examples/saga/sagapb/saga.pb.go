// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: saga.proto

package sagapb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ArgumentDefinition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Types that are assignable to Value:
	//
	//	*ArgumentDefinition_StaticValue
	//	*ArgumentDefinition_ResultReference
	Value isArgumentDefinition_Value `protobuf_oneof:"value"`
}

func (x *ArgumentDefinition) Reset() {
	*x = ArgumentDefinition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArgumentDefinition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArgumentDefinition) ProtoMessage() {}

func (x *ArgumentDefinition) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArgumentDefinition.ProtoReflect.Descriptor instead.
func (*ArgumentDefinition) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{0}
}

func (x *ArgumentDefinition) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (m *ArgumentDefinition) GetValue() isArgumentDefinition_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *ArgumentDefinition) GetStaticValue() *structpb.Value {
	if x, ok := x.GetValue().(*ArgumentDefinition_StaticValue); ok {
		return x.StaticValue
	}
	return nil
}

func (x *ArgumentDefinition) GetResultReference() string {
	if x, ok := x.GetValue().(*ArgumentDefinition_ResultReference); ok {
		return x.ResultReference
	}
	return ""
}

type isArgumentDefinition_Value interface {
	isArgumentDefinition_Value()
}

type ArgumentDefinition_StaticValue struct {
	StaticValue *structpb.Value `protobuf:"bytes,2,opt,name=static_value,json=staticValue,proto3,oneof"`
}

type ArgumentDefinition_ResultReference struct {
	ResultReference string `protobuf:"bytes,3,opt,name=result_reference,json=resultReference,proto3,oneof"`
}

func (*ArgumentDefinition_StaticValue) isArgumentDefinition_Value() {}

func (*ArgumentDefinition_ResultReference) isArgumentDefinition_Value() {}

type TaskDefinition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                    string                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Dependencies          []string              `protobuf:"bytes,2,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
	Method                string                `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
	Arguments             []*ArgumentDefinition `protobuf:"bytes,4,rep,name=arguments,proto3" json:"arguments,omitempty"`
	CompensationMethod    string                `protobuf:"bytes,5,opt,name=compensation_method,json=compensationMethod,proto3" json:"compensation_method,omitempty"`
	CompensationArguments []*ArgumentDefinition `protobuf:"bytes,6,rep,name=compensation_arguments,json=compensationArguments,proto3" json:"compensation_arguments,omitempty"`
}

func (x *TaskDefinition) Reset() {
	*x = TaskDefinition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskDefinition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskDefinition) ProtoMessage() {}

func (x *TaskDefinition) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskDefinition.ProtoReflect.Descriptor instead.
func (*TaskDefinition) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{1}
}

func (x *TaskDefinition) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskDefinition) GetDependencies() []string {
	if x != nil {
		return x.Dependencies
	}
	return nil
}

func (x *TaskDefinition) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *TaskDefinition) GetArguments() []*ArgumentDefinition {
	if x != nil {
		return x.Arguments
	}
	return nil
}

func (x *TaskDefinition) GetCompensationMethod() string {
	if x != nil {
		return x.CompensationMethod
	}
	return ""
}

func (x *TaskDefinition) GetCompensationArguments() []*ArgumentDefinition {
	if x != nil {
		return x.CompensationArguments
	}
	return nil
}

type TaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaId        string           `protobuf:"bytes,1,opt,name=saga_id,json=sagaId,proto3" json:"saga_id,omitempty"`
	TaskId        string           `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	TaskMethod    string           `protobuf:"bytes,3,opt,name=task_method,json=taskMethod,proto3" json:"task_method,omitempty"`
	TaskArguments *structpb.Struct `protobuf:"bytes,4,opt,name=task_arguments,json=taskArguments,proto3" json:"task_arguments,omitempty"`
}

func (x *TaskRequest) Reset() {
	*x = TaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskRequest) ProtoMessage() {}

func (x *TaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskRequest.ProtoReflect.Descriptor instead.
func (*TaskRequest) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{2}
}

func (x *TaskRequest) GetSagaId() string {
	if x != nil {
		return x.SagaId
	}
	return ""
}

func (x *TaskRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *TaskRequest) GetTaskMethod() string {
	if x != nil {
		return x.TaskMethod
	}
	return ""
}

func (x *TaskRequest) GetTaskArguments() *structpb.Struct {
	if x != nil {
		return x.TaskArguments
	}
	return nil
}

type TaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaId          string          `protobuf:"bytes,1,opt,name=saga_id,json=sagaId,proto3" json:"saga_id,omitempty"`
	TaskId          string          `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	TaskAborted     bool            `protobuf:"varint,3,opt,name=task_aborted,json=taskAborted,proto3" json:"task_aborted,omitempty"`
	TaskResult      *structpb.Value `protobuf:"bytes,4,opt,name=task_result,json=taskResult,proto3" json:"task_result,omitempty"`
	TaskAbortReason *structpb.Value `protobuf:"bytes,5,opt,name=task_abort_reason,json=taskAbortReason,proto3" json:"task_abort_reason,omitempty"`
}

func (x *TaskResponse) Reset() {
	*x = TaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResponse) ProtoMessage() {}

func (x *TaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResponse.ProtoReflect.Descriptor instead.
func (*TaskResponse) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{3}
}

func (x *TaskResponse) GetSagaId() string {
	if x != nil {
		return x.SagaId
	}
	return ""
}

func (x *TaskResponse) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *TaskResponse) GetTaskAborted() bool {
	if x != nil {
		return x.TaskAborted
	}
	return false
}

func (x *TaskResponse) GetTaskResult() *structpb.Value {
	if x != nil {
		return x.TaskResult
	}
	return nil
}

func (x *TaskResponse) GetTaskAbortReason() *structpb.Value {
	if x != nil {
		return x.TaskAbortReason
	}
	return nil
}

type CompensationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaId                    string           `protobuf:"bytes,1,opt,name=saga_id,json=sagaId,proto3" json:"saga_id,omitempty"`
	TaskId                    string           `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	TaskArguments             *structpb.Struct `protobuf:"bytes,3,opt,name=task_arguments,json=taskArguments,proto3" json:"task_arguments,omitempty"`
	TaskResult                *structpb.Value  `protobuf:"bytes,4,opt,name=task_result,json=taskResult,proto3" json:"task_result,omitempty"`
	TaskCompensationMethod    string           `protobuf:"bytes,5,opt,name=task_compensation_method,json=taskCompensationMethod,proto3" json:"task_compensation_method,omitempty"`
	TaskCompensationArguments *structpb.Struct `protobuf:"bytes,6,opt,name=task_compensation_arguments,json=taskCompensationArguments,proto3" json:"task_compensation_arguments,omitempty"`
}

func (x *CompensationRequest) Reset() {
	*x = CompensationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompensationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompensationRequest) ProtoMessage() {}

func (x *CompensationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompensationRequest.ProtoReflect.Descriptor instead.
func (*CompensationRequest) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{4}
}

func (x *CompensationRequest) GetSagaId() string {
	if x != nil {
		return x.SagaId
	}
	return ""
}

func (x *CompensationRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *CompensationRequest) GetTaskArguments() *structpb.Struct {
	if x != nil {
		return x.TaskArguments
	}
	return nil
}

func (x *CompensationRequest) GetTaskResult() *structpb.Value {
	if x != nil {
		return x.TaskResult
	}
	return nil
}

func (x *CompensationRequest) GetTaskCompensationMethod() string {
	if x != nil {
		return x.TaskCompensationMethod
	}
	return ""
}

func (x *CompensationRequest) GetTaskCompensationArguments() *structpb.Struct {
	if x != nil {
		return x.TaskCompensationArguments
	}
	return nil
}

type CompensationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaId string `protobuf:"bytes,1,opt,name=saga_id,json=sagaId,proto3" json:"saga_id,omitempty"`
	TaskId string `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (x *CompensationResponse) Reset() {
	*x = CompensationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompensationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompensationResponse) ProtoMessage() {}

func (x *CompensationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompensationResponse.ProtoReflect.Descriptor instead.
func (*CompensationResponse) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{5}
}

func (x *CompensationResponse) GetSagaId() string {
	if x != nil {
		return x.SagaId
	}
	return ""
}

func (x *CompensationResponse) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

type SagaBegun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskDefinitions []*TaskDefinition `protobuf:"bytes,2,rep,name=task_definitions,json=taskDefinitions,proto3" json:"task_definitions,omitempty"`
}

func (x *SagaBegun) Reset() {
	*x = SagaBegun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SagaBegun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SagaBegun) ProtoMessage() {}

func (x *SagaBegun) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SagaBegun.ProtoReflect.Descriptor instead.
func (*SagaBegun) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{6}
}

func (x *SagaBegun) GetTaskDefinitions() []*TaskDefinition {
	if x != nil {
		return x.TaskDefinitions
	}
	return nil
}

type SagaEnded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SagaEnded) Reset() {
	*x = SagaEnded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SagaEnded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SagaEnded) ProtoMessage() {}

func (x *SagaEnded) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SagaEnded.ProtoReflect.Descriptor instead.
func (*SagaEnded) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{7}
}

type TaskBegun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *TaskBegun) Reset() {
	*x = TaskBegun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskBegun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskBegun) ProtoMessage() {}

func (x *TaskBegun) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskBegun.ProtoReflect.Descriptor instead.
func (*TaskBegun) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{8}
}

func (x *TaskBegun) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type TaskEnded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Result *structpb.Value `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *TaskEnded) Reset() {
	*x = TaskEnded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskEnded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskEnded) ProtoMessage() {}

func (x *TaskEnded) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskEnded.ProtoReflect.Descriptor instead.
func (*TaskEnded) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{9}
}

func (x *TaskEnded) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskEnded) GetResult() *structpb.Value {
	if x != nil {
		return x.Result
	}
	return nil
}

type TaskAborted struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Reason *structpb.Value `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *TaskAborted) Reset() {
	*x = TaskAborted{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskAborted) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskAborted) ProtoMessage() {}

func (x *TaskAborted) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskAborted.ProtoReflect.Descriptor instead.
func (*TaskAborted) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{10}
}

func (x *TaskAborted) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskAborted) GetReason() *structpb.Value {
	if x != nil {
		return x.Reason
	}
	return nil
}

type CompensationBegun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CompensationBegun) Reset() {
	*x = CompensationBegun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompensationBegun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompensationBegun) ProtoMessage() {}

func (x *CompensationBegun) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompensationBegun.ProtoReflect.Descriptor instead.
func (*CompensationBegun) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{11}
}

func (x *CompensationBegun) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CompensationEnded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CompensationEnded) Reset() {
	*x = CompensationEnded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saga_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompensationEnded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompensationEnded) ProtoMessage() {}

func (x *CompensationEnded) ProtoReflect() protoreflect.Message {
	mi := &file_saga_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompensationEnded.ProtoReflect.Descriptor instead.
func (*CompensationEnded) Descriptor() ([]byte, []int) {
	return file_saga_proto_rawDescGZIP(), []int{12}
}

func (x *CompensationEnded) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_saga_proto protoreflect.FileDescriptor

var file_saga_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x73, 0x61, 0x67, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x24, 0x72, 0x6e,
	0x6f, 0x76, 0x61, 0x74, 0x6f, 0x72, 0x6f, 0x76, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x73, 0x61,
	0x67, 0x61, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x9b, 0x01, 0x0a, 0x12, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x66,
	0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x0c, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2b, 0x0a, 0x10, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x66, 0x65,
	0x72, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xd6,
	0x02, 0x0a, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65,
	0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x56, 0x0a,
	0x09, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x38, 0x2e, 0x72, 0x6e, 0x6f, 0x76, 0x61, 0x74, 0x6f, 0x72, 0x6f, 0x76, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x73, 0x2e, 0x73, 0x61, 0x67, 0x61, 0x2e, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x61, 0x72, 0x67, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2f, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x12, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x6f, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x6e,
	0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x72, 0x6e, 0x6f, 0x76, 0x61, 0x74, 0x6f,
	0x72, 0x6f, 0x76, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x73, 0x61, 0x67, 0x61, 0x2e, 0x41, 0x72,
	0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x15, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x72,
	0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0xa0, 0x01, 0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x61, 0x67, 0x61, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x61, 0x67, 0x61, 0x49, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x73,
	0x6b, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x74, 0x61, 0x73, 0x6b, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x3e, 0x0a, 0x0e, 0x74, 0x61,
	0x73, 0x6b, 0x5f, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x0d, 0x74, 0x61, 0x73,
	0x6b, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0xe0, 0x01, 0x0a, 0x0c, 0x54,
	0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x73,
	0x61, 0x67, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x61,
	0x67, 0x61, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x21, 0x0a,
	0x0c, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x61, 0x62, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x74, 0x61, 0x73, 0x6b, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x65, 0x64,
	0x12, 0x37, 0x0a, 0x0b, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x74,
	0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x42, 0x0a, 0x11, 0x74, 0x61, 0x73,
	0x6b, 0x5f, 0x61, 0x62, 0x6f, 0x72, 0x74, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0f, 0x74, 0x61,
	0x73, 0x6b, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0xd3, 0x02,
	0x0a, 0x13, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x61, 0x67, 0x61, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x61, 0x67, 0x61, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x3e, 0x0a, 0x0e, 0x74, 0x61, 0x73, 0x6b, 0x5f,
	0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x0d, 0x74, 0x61, 0x73, 0x6b, 0x41, 0x72,
	0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x37, 0x0a, 0x0b, 0x74, 0x61, 0x73, 0x6b, 0x5f,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x38, 0x0a, 0x18, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x16, 0x74, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x57, 0x0a, 0x1b, 0x74, 0x61,
	0x73, 0x6b, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x19, 0x74, 0x61, 0x73, 0x6b, 0x43, 0x6f,
	0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x22, 0x48, 0x0a, 0x14, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x73,
	0x61, 0x67, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x61,
	0x67, 0x61, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x22, 0x6c, 0x0a,
	0x09, 0x53, 0x61, 0x67, 0x61, 0x42, 0x65, 0x67, 0x75, 0x6e, 0x12, 0x5f, 0x0a, 0x10, 0x74, 0x61,
	0x73, 0x6b, 0x5f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x72, 0x6e, 0x6f, 0x76, 0x61, 0x74, 0x6f, 0x72, 0x6f,
	0x76, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x73, 0x61, 0x67, 0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x74, 0x61, 0x73, 0x6b,
	0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x0b, 0x0a, 0x09, 0x53,
	0x61, 0x67, 0x61, 0x45, 0x6e, 0x64, 0x65, 0x64, 0x22, 0x1b, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b,
	0x42, 0x65, 0x67, 0x75, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4b, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x6e, 0x64,
	0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x22, 0x4d, 0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x2e, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x22, 0x23, 0x0a, 0x11, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x65, 0x67, 0x75, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x23, 0x0a, 0x11, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e,
	0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_saga_proto_rawDescOnce sync.Once
	file_saga_proto_rawDescData = file_saga_proto_rawDesc
)

func file_saga_proto_rawDescGZIP() []byte {
	file_saga_proto_rawDescOnce.Do(func() {
		file_saga_proto_rawDescData = protoimpl.X.CompressGZIP(file_saga_proto_rawDescData)
	})
	return file_saga_proto_rawDescData
}

var file_saga_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_saga_proto_goTypes = []any{
	(*ArgumentDefinition)(nil),   // 0: rnovatorov.eventsource.examples.saga.ArgumentDefinition
	(*TaskDefinition)(nil),       // 1: rnovatorov.eventsource.examples.saga.TaskDefinition
	(*TaskRequest)(nil),          // 2: rnovatorov.eventsource.examples.saga.TaskRequest
	(*TaskResponse)(nil),         // 3: rnovatorov.eventsource.examples.saga.TaskResponse
	(*CompensationRequest)(nil),  // 4: rnovatorov.eventsource.examples.saga.CompensationRequest
	(*CompensationResponse)(nil), // 5: rnovatorov.eventsource.examples.saga.CompensationResponse
	(*SagaBegun)(nil),            // 6: rnovatorov.eventsource.examples.saga.SagaBegun
	(*SagaEnded)(nil),            // 7: rnovatorov.eventsource.examples.saga.SagaEnded
	(*TaskBegun)(nil),            // 8: rnovatorov.eventsource.examples.saga.TaskBegun
	(*TaskEnded)(nil),            // 9: rnovatorov.eventsource.examples.saga.TaskEnded
	(*TaskAborted)(nil),          // 10: rnovatorov.eventsource.examples.saga.TaskAborted
	(*CompensationBegun)(nil),    // 11: rnovatorov.eventsource.examples.saga.CompensationBegun
	(*CompensationEnded)(nil),    // 12: rnovatorov.eventsource.examples.saga.CompensationEnded
	(*structpb.Value)(nil),       // 13: google.protobuf.Value
	(*structpb.Struct)(nil),      // 14: google.protobuf.Struct
}
var file_saga_proto_depIdxs = []int32{
	13, // 0: rnovatorov.eventsource.examples.saga.ArgumentDefinition.static_value:type_name -> google.protobuf.Value
	0,  // 1: rnovatorov.eventsource.examples.saga.TaskDefinition.arguments:type_name -> rnovatorov.eventsource.examples.saga.ArgumentDefinition
	0,  // 2: rnovatorov.eventsource.examples.saga.TaskDefinition.compensation_arguments:type_name -> rnovatorov.eventsource.examples.saga.ArgumentDefinition
	14, // 3: rnovatorov.eventsource.examples.saga.TaskRequest.task_arguments:type_name -> google.protobuf.Struct
	13, // 4: rnovatorov.eventsource.examples.saga.TaskResponse.task_result:type_name -> google.protobuf.Value
	13, // 5: rnovatorov.eventsource.examples.saga.TaskResponse.task_abort_reason:type_name -> google.protobuf.Value
	14, // 6: rnovatorov.eventsource.examples.saga.CompensationRequest.task_arguments:type_name -> google.protobuf.Struct
	13, // 7: rnovatorov.eventsource.examples.saga.CompensationRequest.task_result:type_name -> google.protobuf.Value
	14, // 8: rnovatorov.eventsource.examples.saga.CompensationRequest.task_compensation_arguments:type_name -> google.protobuf.Struct
	1,  // 9: rnovatorov.eventsource.examples.saga.SagaBegun.task_definitions:type_name -> rnovatorov.eventsource.examples.saga.TaskDefinition
	13, // 10: rnovatorov.eventsource.examples.saga.TaskEnded.result:type_name -> google.protobuf.Value
	13, // 11: rnovatorov.eventsource.examples.saga.TaskAborted.reason:type_name -> google.protobuf.Value
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_saga_proto_init() }
func file_saga_proto_init() {
	if File_saga_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_saga_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ArgumentDefinition); i {
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
		file_saga_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TaskDefinition); i {
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
		file_saga_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*TaskRequest); i {
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
		file_saga_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*TaskResponse); i {
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
		file_saga_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*CompensationRequest); i {
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
		file_saga_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*CompensationResponse); i {
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
		file_saga_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*SagaBegun); i {
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
		file_saga_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SagaEnded); i {
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
		file_saga_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*TaskBegun); i {
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
		file_saga_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*TaskEnded); i {
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
		file_saga_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*TaskAborted); i {
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
		file_saga_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*CompensationBegun); i {
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
		file_saga_proto_msgTypes[12].Exporter = func(v any, i int) any {
			switch v := v.(*CompensationEnded); i {
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
	file_saga_proto_msgTypes[0].OneofWrappers = []any{
		(*ArgumentDefinition_StaticValue)(nil),
		(*ArgumentDefinition_ResultReference)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_saga_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_saga_proto_goTypes,
		DependencyIndexes: file_saga_proto_depIdxs,
		MessageInfos:      file_saga_proto_msgTypes,
	}.Build()
	File_saga_proto = out.File
	file_saga_proto_rawDesc = nil
	file_saga_proto_goTypes = nil
	file_saga_proto_depIdxs = nil
}
