// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.0
// 	protoc        v5.26.1
// source: internal/mygrpc/proto/mygrpc.proto

package proto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metrica_MType int32

const (
	Metrica_COUNTER Metrica_MType = 0
	Metrica_GAUGE   Metrica_MType = 1
)

// Enum value maps for Metrica_MType.
var (
	Metrica_MType_name = map[int32]string{
		0: "COUNTER",
		1: "GAUGE",
	}
	Metrica_MType_value = map[string]int32{
		"COUNTER": 0,
		"GAUGE":   1,
	}
)

func (x Metrica_MType) Enum() *Metrica_MType {
	p := new(Metrica_MType)
	*p = x
	return p
}

func (x Metrica_MType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Metrica_MType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_mygrpc_proto_mygrpc_proto_enumTypes[0].Descriptor()
}

func (Metrica_MType) Type() protoreflect.EnumType {
	return &file_internal_mygrpc_proto_mygrpc_proto_enumTypes[0]
}

func (x Metrica_MType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Metrica_MType.Descriptor instead.
func (Metrica_MType) EnumDescriptor() ([]byte, []int) {
	return file_internal_mygrpc_proto_mygrpc_proto_rawDescGZIP(), []int{0, 0}
}

type Metrica struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Mtype Metrica_MType `protobuf:"varint,2,opt,name=mtype,proto3,enum=mygrpc.Metrica_MType" json:"mtype,omitempty"`
	Delta int64         `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`
	Value float64       `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Metrica) Reset() {
	*x = Metrica{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_mygrpc_proto_mygrpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metrica) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metrica) ProtoMessage() {}

func (x *Metrica) ProtoReflect() protoreflect.Message {
	mi := &file_internal_mygrpc_proto_mygrpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metrica.ProtoReflect.Descriptor instead.
func (*Metrica) Descriptor() ([]byte, []int) {
	return file_internal_mygrpc_proto_mygrpc_proto_rawDescGZIP(), []int{0}
}

func (x *Metrica) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metrica) GetMtype() Metrica_MType {
	if x != nil {
		return x.Mtype
	}
	return Metrica_COUNTER
}

func (x *Metrica) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metrica) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metrica `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_mygrpc_proto_mygrpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_mygrpc_proto_mygrpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_internal_mygrpc_proto_mygrpc_proto_rawDescGZIP(), []int{1}
}

func (x *Request) GetMetrics() []*Metrica {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type UpdatesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *UpdatesResponse) Reset() {
	*x = UpdatesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_mygrpc_proto_mygrpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatesResponse) ProtoMessage() {}

func (x *UpdatesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_mygrpc_proto_mygrpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatesResponse.ProtoReflect.Descriptor instead.
func (*UpdatesResponse) Descriptor() ([]byte, []int) {
	return file_internal_mygrpc_proto_mygrpc_proto_rawDescGZIP(), []int{2}
}

func (x *UpdatesResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_internal_mygrpc_proto_mygrpc_proto protoreflect.FileDescriptor

var file_internal_mygrpc_proto_mygrpc_proto_rawDesc = []byte{
	0x0a, 0x22, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x79, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x79, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6d, 0x79, 0x67, 0x72, 0x70, 0x63, 0x22, 0x93, 0x01, 0x0a,
	0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2b, 0x0a, 0x05, 0x6d, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x6d, 0x79, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x61, 0x2e, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05,
	0x6d, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x1f, 0x0a, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f,
	0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x41, 0x55, 0x47, 0x45,
	0x10, 0x01, 0x22, 0x34, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a,
	0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x6d, 0x79, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x61, 0x52,
	0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x25, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32,
	0x45, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x3a, 0x0a, 0x0e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x73, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x0f, 0x2e, 0x6d,
	0x79, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x6d, 0x79, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x17, 0x5a, 0x15, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x6d, 0x79, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_mygrpc_proto_mygrpc_proto_rawDescOnce sync.Once
	file_internal_mygrpc_proto_mygrpc_proto_rawDescData = file_internal_mygrpc_proto_mygrpc_proto_rawDesc
)

func file_internal_mygrpc_proto_mygrpc_proto_rawDescGZIP() []byte {
	file_internal_mygrpc_proto_mygrpc_proto_rawDescOnce.Do(func() {
		file_internal_mygrpc_proto_mygrpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_mygrpc_proto_mygrpc_proto_rawDescData)
	})
	return file_internal_mygrpc_proto_mygrpc_proto_rawDescData
}

var file_internal_mygrpc_proto_mygrpc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_mygrpc_proto_mygrpc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_mygrpc_proto_mygrpc_proto_goTypes = []interface{}{
	(Metrica_MType)(0),      // 0: mygrpc.Metrica.MType
	(*Metrica)(nil),         // 1: mygrpc.Metrica
	(*Request)(nil),         // 2: mygrpc.Request
	(*UpdatesResponse)(nil), // 3: mygrpc.UpdatesResponse
}
var file_internal_mygrpc_proto_mygrpc_proto_depIdxs = []int32{
	0, // 0: mygrpc.Metrica.mtype:type_name -> mygrpc.Metrica.MType
	1, // 1: mygrpc.Request.metrics:type_name -> mygrpc.Metrica
	2, // 2: mygrpc.Metrics.UpdatesMetrics:input_type -> mygrpc.Request
	3, // 3: mygrpc.Metrics.UpdatesMetrics:output_type -> mygrpc.UpdatesResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_mygrpc_proto_mygrpc_proto_init() }
func file_internal_mygrpc_proto_mygrpc_proto_init() {
	if File_internal_mygrpc_proto_mygrpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_mygrpc_proto_mygrpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metrica); i {
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
		file_internal_mygrpc_proto_mygrpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_internal_mygrpc_proto_mygrpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatesResponse); i {
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
			RawDescriptor: file_internal_mygrpc_proto_mygrpc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_mygrpc_proto_mygrpc_proto_goTypes,
		DependencyIndexes: file_internal_mygrpc_proto_mygrpc_proto_depIdxs,
		EnumInfos:         file_internal_mygrpc_proto_mygrpc_proto_enumTypes,
		MessageInfos:      file_internal_mygrpc_proto_mygrpc_proto_msgTypes,
	}.Build()
	File_internal_mygrpc_proto_mygrpc_proto = out.File
	file_internal_mygrpc_proto_mygrpc_proto_rawDesc = nil
	file_internal_mygrpc_proto_mygrpc_proto_goTypes = nil
	file_internal_mygrpc_proto_mygrpc_proto_depIdxs = nil
}
