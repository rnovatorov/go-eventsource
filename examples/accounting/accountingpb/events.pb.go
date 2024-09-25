// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: events.proto

package accountingpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Book_AccountType int32

const (
	Book_ACCOUNT_TYPE_UNKNOWN   Book_AccountType = 0
	Book_ACCOUNT_TYPE_CAPITAL   Book_AccountType = 1
	Book_ACCOUNT_TYPE_ASSET     Book_AccountType = 2
	Book_ACCOUNT_TYPE_LIABILITY Book_AccountType = 3
	Book_ACCOUNT_TYPE_INCOME    Book_AccountType = 4
	Book_ACCOUNT_TYPE_EXPENSE   Book_AccountType = 5
)

// Enum value maps for Book_AccountType.
var (
	Book_AccountType_name = map[int32]string{
		0: "ACCOUNT_TYPE_UNKNOWN",
		1: "ACCOUNT_TYPE_CAPITAL",
		2: "ACCOUNT_TYPE_ASSET",
		3: "ACCOUNT_TYPE_LIABILITY",
		4: "ACCOUNT_TYPE_INCOME",
		5: "ACCOUNT_TYPE_EXPENSE",
	}
	Book_AccountType_value = map[string]int32{
		"ACCOUNT_TYPE_UNKNOWN":   0,
		"ACCOUNT_TYPE_CAPITAL":   1,
		"ACCOUNT_TYPE_ASSET":     2,
		"ACCOUNT_TYPE_LIABILITY": 3,
		"ACCOUNT_TYPE_INCOME":    4,
		"ACCOUNT_TYPE_EXPENSE":   5,
	}
)

func (x Book_AccountType) Enum() *Book_AccountType {
	p := new(Book_AccountType)
	*p = x
	return p
}

func (x Book_AccountType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Book_AccountType) Descriptor() protoreflect.EnumDescriptor {
	return file_events_proto_enumTypes[0].Descriptor()
}

func (Book_AccountType) Type() protoreflect.EnumType {
	return &file_events_proto_enumTypes[0]
}

func (x Book_AccountType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Book_AccountType.Descriptor instead.
func (Book_AccountType) EnumDescriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0, 0}
}

type Book struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Book) Reset() {
	*x = Book{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book) ProtoMessage() {}

func (x *Book) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book.ProtoReflect.Descriptor instead.
func (*Book) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0}
}

type Book_Created struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Book_Created) Reset() {
	*x = Book_Created{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book_Created) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book_Created) ProtoMessage() {}

func (x *Book_Created) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book_Created.ProtoReflect.Descriptor instead.
func (*Book_Created) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Book_Created) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type Book_Closed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Book_Closed) Reset() {
	*x = Book_Closed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book_Closed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book_Closed) ProtoMessage() {}

func (x *Book_Closed) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book_Closed.ProtoReflect.Descriptor instead.
func (*Book_Closed) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0, 1}
}

type Book_AccountAdded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type Book_AccountType `protobuf:"varint,2,opt,name=type,proto3,enum=rnovatorov.eventsource.accounting.Book_AccountType" json:"type,omitempty"`
}

func (x *Book_AccountAdded) Reset() {
	*x = Book_AccountAdded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book_AccountAdded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book_AccountAdded) ProtoMessage() {}

func (x *Book_AccountAdded) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book_AccountAdded.ProtoReflect.Descriptor instead.
func (*Book_AccountAdded) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Book_AccountAdded) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Book_AccountAdded) GetType() Book_AccountType {
	if x != nil {
		return x.Type
	}
	return Book_ACCOUNT_TYPE_UNKNOWN
}

type Book_TransactionEntered struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp       *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	AccountDebited  string                 `protobuf:"bytes,2,opt,name=account_debited,json=accountDebited,proto3" json:"account_debited,omitempty"`
	AccountCredited string                 `protobuf:"bytes,3,opt,name=account_credited,json=accountCredited,proto3" json:"account_credited,omitempty"`
	Amount          uint64                 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Book_TransactionEntered) Reset() {
	*x = Book_TransactionEntered{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book_TransactionEntered) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book_TransactionEntered) ProtoMessage() {}

func (x *Book_TransactionEntered) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book_TransactionEntered.ProtoReflect.Descriptor instead.
func (*Book_TransactionEntered) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0, 3}
}

func (x *Book_TransactionEntered) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Book_TransactionEntered) GetAccountDebited() string {
	if x != nil {
		return x.AccountDebited
	}
	return ""
}

func (x *Book_TransactionEntered) GetAccountCredited() string {
	if x != nil {
		return x.AccountCredited
	}
	return ""
}

func (x *Book_TransactionEntered) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

var File_events_proto protoreflect.FileDescriptor

var file_events_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x21,
	0x72, 0x6e, 0x6f, 0x76, 0x61, 0x74, 0x6f, 0x72, 0x6f, 0x76, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e,
	0x67, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x92, 0x04, 0x0a, 0x04, 0x42, 0x6f, 0x6f, 0x6b, 0x1a, 0x2b, 0x0a, 0x07, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x08, 0x0a, 0x06, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x64, 0x1a, 0x6b, 0x0a, 0x0c, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x64, 0x64,
	0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x47, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x33, 0x2e, 0x72, 0x6e, 0x6f, 0x76, 0x61, 0x74, 0x6f, 0x72, 0x6f,
	0x76, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x42, 0x6f, 0x6f, 0x6b, 0x2e, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x1a,
	0xba, 0x01, 0x0a, 0x12, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45,
	0x6e, 0x74, 0x65, 0x72, 0x65, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x27, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x64, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x44, 0x65, 0x62, 0x69, 0x74, 0x65, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x64,
	0x69, 0x74, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xa8, 0x01, 0x0a,
	0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x14,
	0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e,
	0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x41, 0x50, 0x49, 0x54, 0x41, 0x4c, 0x10, 0x01,
	0x12, 0x16, 0x0a, 0x12, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x41, 0x53, 0x53, 0x45, 0x54, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x43, 0x43, 0x4f,
	0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x49, 0x41, 0x42, 0x49, 0x4c, 0x49,
	0x54, 0x59, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x43, 0x4f, 0x4d, 0x45, 0x10, 0x04, 0x12, 0x18, 0x0a,
	0x14, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45, 0x58,
	0x50, 0x45, 0x4e, 0x53, 0x45, 0x10, 0x05, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x3b, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_events_proto_rawDescOnce sync.Once
	file_events_proto_rawDescData = file_events_proto_rawDesc
)

func file_events_proto_rawDescGZIP() []byte {
	file_events_proto_rawDescOnce.Do(func() {
		file_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_events_proto_rawDescData)
	})
	return file_events_proto_rawDescData
}

var file_events_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_events_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_events_proto_goTypes = []any{
	(Book_AccountType)(0),           // 0: rnovatorov.eventsource.accounting.Book.AccountType
	(*Book)(nil),                    // 1: rnovatorov.eventsource.accounting.Book
	(*Book_Created)(nil),            // 2: rnovatorov.eventsource.accounting.Book.Created
	(*Book_Closed)(nil),             // 3: rnovatorov.eventsource.accounting.Book.Closed
	(*Book_AccountAdded)(nil),       // 4: rnovatorov.eventsource.accounting.Book.AccountAdded
	(*Book_TransactionEntered)(nil), // 5: rnovatorov.eventsource.accounting.Book.TransactionEntered
	(*timestamppb.Timestamp)(nil),   // 6: google.protobuf.Timestamp
}
var file_events_proto_depIdxs = []int32{
	0, // 0: rnovatorov.eventsource.accounting.Book.AccountAdded.type:type_name -> rnovatorov.eventsource.accounting.Book.AccountType
	6, // 1: rnovatorov.eventsource.accounting.Book.TransactionEntered.timestamp:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_events_proto_init() }
func file_events_proto_init() {
	if File_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_events_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Book); i {
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
		file_events_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Book_Created); i {
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
		file_events_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Book_Closed); i {
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
		file_events_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Book_AccountAdded); i {
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
		file_events_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Book_TransactionEntered); i {
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
			RawDescriptor: file_events_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_proto_goTypes,
		DependencyIndexes: file_events_proto_depIdxs,
		EnumInfos:         file_events_proto_enumTypes,
		MessageInfos:      file_events_proto_msgTypes,
	}.Build()
	File_events_proto = out.File
	file_events_proto_rawDesc = nil
	file_events_proto_goTypes = nil
	file_events_proto_depIdxs = nil
}