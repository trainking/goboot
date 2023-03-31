// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: gameserver.proto

package pb

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

// OpCode 操作符定义
type OpCode int32

const (
	OpCode_None OpCode = 0
	OpCode_Ping OpCode = 1
	OpCode_Pong OpCode = 2
)

// Enum value maps for OpCode.
var (
	OpCode_name = map[int32]string{
		0: "None",
		1: "Ping",
		2: "Pong",
	}
	OpCode_value = map[string]int32{
		"None": 0,
		"Ping": 1,
		"Pong": 2,
	}
)

func (x OpCode) Enum() *OpCode {
	p := new(OpCode)
	*p = x
	return p
}

func (x OpCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpCode) Descriptor() protoreflect.EnumDescriptor {
	return file_gameserver_proto_enumTypes[0].Descriptor()
}

func (OpCode) Type() protoreflect.EnumType {
	return &file_gameserver_proto_enumTypes[0]
}

func (x OpCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpCode.Descriptor instead.
func (OpCode) EnumDescriptor() ([]byte, []int) {
	return file_gameserver_proto_rawDescGZIP(), []int{0}
}

// C2S_Ping 客户端向服务器发送心跳
type C2S_Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TickTime int64 `protobuf:"varint,1,opt,name=TickTime,proto3" json:"TickTime"`
}

func (x *C2S_Ping) Reset() {
	*x = C2S_Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gameserver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2S_Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2S_Ping) ProtoMessage() {}

func (x *C2S_Ping) ProtoReflect() protoreflect.Message {
	mi := &file_gameserver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2S_Ping.ProtoReflect.Descriptor instead.
func (*C2S_Ping) Descriptor() ([]byte, []int) {
	return file_gameserver_proto_rawDescGZIP(), []int{0}
}

func (x *C2S_Ping) GetTickTime() int64 {
	if x != nil {
		return x.TickTime
	}
	return 0
}

type S2C_Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OK bool `protobuf:"varint,1,opt,name=OK,proto3" json:"OK"`
}

func (x *S2C_Pong) Reset() {
	*x = S2C_Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gameserver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2C_Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2C_Pong) ProtoMessage() {}

func (x *S2C_Pong) ProtoReflect() protoreflect.Message {
	mi := &file_gameserver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2C_Pong.ProtoReflect.Descriptor instead.
func (*S2C_Pong) Descriptor() ([]byte, []int) {
	return file_gameserver_proto_rawDescGZIP(), []int{1}
}

func (x *S2C_Pong) GetOK() bool {
	if x != nil {
		return x.OK
	}
	return false
}

var File_gameserver_proto protoreflect.FileDescriptor

var file_gameserver_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x26, 0x0a, 0x08, 0x43, 0x32, 0x53, 0x5f, 0x50, 0x69,
	0x6e, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x54, 0x69, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x54, 0x69, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x1a,
	0x0a, 0x08, 0x53, 0x32, 0x43, 0x5f, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x4f, 0x4b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x4b, 0x2a, 0x26, 0x0a, 0x06, 0x4f, 0x70,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67,
	0x10, 0x02, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_gameserver_proto_rawDescOnce sync.Once
	file_gameserver_proto_rawDescData = file_gameserver_proto_rawDesc
)

func file_gameserver_proto_rawDescGZIP() []byte {
	file_gameserver_proto_rawDescOnce.Do(func() {
		file_gameserver_proto_rawDescData = protoimpl.X.CompressGZIP(file_gameserver_proto_rawDescData)
	})
	return file_gameserver_proto_rawDescData
}

var file_gameserver_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gameserver_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gameserver_proto_goTypes = []interface{}{
	(OpCode)(0),      // 0: pb.OpCode
	(*C2S_Ping)(nil), // 1: pb.C2S_Ping
	(*S2C_Pong)(nil), // 2: pb.S2C_Pong
}
var file_gameserver_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gameserver_proto_init() }
func file_gameserver_proto_init() {
	if File_gameserver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gameserver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2S_Ping); i {
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
		file_gameserver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2C_Pong); i {
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
			RawDescriptor: file_gameserver_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gameserver_proto_goTypes,
		DependencyIndexes: file_gameserver_proto_depIdxs,
		EnumInfos:         file_gameserver_proto_enumTypes,
		MessageInfos:      file_gameserver_proto_msgTypes,
	}.Build()
	File_gameserver_proto = out.File
	file_gameserver_proto_rawDesc = nil
	file_gameserver_proto_goTypes = nil
	file_gameserver_proto_depIdxs = nil
}
