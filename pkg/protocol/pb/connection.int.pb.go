// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.2
// source: connection.int.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 消息请求
type TransferMessageReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DeviceId      int64                  `protobuf:"varint,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"` //设备id
	Message       *Message               `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`                    //数据
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransferMessageReq) Reset() {
	*x = TransferMessageReq{}
	mi := &file_connection_int_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransferMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferMessageReq) ProtoMessage() {}

func (x *TransferMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_connection_int_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferMessageReq.ProtoReflect.Descriptor instead.
func (*TransferMessageReq) Descriptor() ([]byte, []int) {
	return file_connection_int_proto_rawDescGZIP(), []int{0}
}

func (x *TransferMessageReq) GetDeviceId() int64 {
	if x != nil {
		return x.DeviceId
	}
	return 0
}

func (x *TransferMessageReq) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

// 房间推送
type PushRoomMsg struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RootId        int64                  `protobuf:"varint,1,opt,name=root_id,json=rootId,proto3" json:"root_id,omitempty"` //房间id
	Message       *Message               `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`              // 数据
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PushRoomMsg) Reset() {
	*x = PushRoomMsg{}
	mi := &file_connection_int_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PushRoomMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushRoomMsg) ProtoMessage() {}

func (x *PushRoomMsg) ProtoReflect() protoreflect.Message {
	mi := &file_connection_int_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushRoomMsg.ProtoReflect.Descriptor instead.
func (*PushRoomMsg) Descriptor() ([]byte, []int) {
	return file_connection_int_proto_rawDescGZIP(), []int{1}
}

func (x *PushRoomMsg) GetRootId() int64 {
	if x != nil {
		return x.RootId
	}
	return 0
}

func (x *PushRoomMsg) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

// 推送全部
type PushAllMsg struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       *Message               `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"` // 数据
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PushAllMsg) Reset() {
	*x = PushAllMsg{}
	mi := &file_connection_int_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PushAllMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushAllMsg) ProtoMessage() {}

func (x *PushAllMsg) ProtoReflect() protoreflect.Message {
	mi := &file_connection_int_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushAllMsg.ProtoReflect.Descriptor instead.
func (*PushAllMsg) Descriptor() ([]byte, []int) {
	return file_connection_int_proto_rawDescGZIP(), []int{2}
}

func (x *PushAllMsg) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_connection_int_proto protoreflect.FileDescriptor

var file_connection_int_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x69, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58, 0x0a, 0x12, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71,
	0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x4d, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x52, 0x6f, 0x6f, 0x6d,
	0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x74, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x33, 0x0a, 0x0a, 0x50, 0x75, 0x73, 0x68, 0x41, 0x6c, 0x6c, 0x4d, 0x73,
	0x67, 0x12, 0x25, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x52, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x74, 0x12, 0x41, 0x0a, 0x0f, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x66, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x2e, 0x70,
	0x62, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x11, 0x5a, 0x0f,
	0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_connection_int_proto_rawDescOnce sync.Once
	file_connection_int_proto_rawDescData []byte
)

func file_connection_int_proto_rawDescGZIP() []byte {
	file_connection_int_proto_rawDescOnce.Do(func() {
		file_connection_int_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_connection_int_proto_rawDesc), len(file_connection_int_proto_rawDesc)))
	})
	return file_connection_int_proto_rawDescData
}

var file_connection_int_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_connection_int_proto_goTypes = []any{
	(*TransferMessageReq)(nil), // 0: pb.TransferMessageReq
	(*PushRoomMsg)(nil),        // 1: pb.PushRoomMsg
	(*PushAllMsg)(nil),         // 2: pb.PushAllMsg
	(*Message)(nil),            // 3: pb.Message
	(*emptypb.Empty)(nil),      // 4: google.protobuf.Empty
}
var file_connection_int_proto_depIdxs = []int32{
	3, // 0: pb.TransferMessageReq.message:type_name -> pb.Message
	3, // 1: pb.PushRoomMsg.message:type_name -> pb.Message
	3, // 2: pb.PushAllMsg.message:type_name -> pb.Message
	0, // 3: pb.ConnectionInt.TransferMessage:input_type -> pb.TransferMessageReq
	4, // 4: pb.ConnectionInt.TransferMessage:output_type -> google.protobuf.Empty
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_connection_int_proto_init() }
func file_connection_int_proto_init() {
	if File_connection_int_proto != nil {
		return
	}
	file_message_ext_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_connection_int_proto_rawDesc), len(file_connection_int_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_connection_int_proto_goTypes,
		DependencyIndexes: file_connection_int_proto_depIdxs,
		MessageInfos:      file_connection_int_proto_msgTypes,
	}.Build()
	File_connection_int_proto = out.File
	file_connection_int_proto_goTypes = nil
	file_connection_int_proto_depIdxs = nil
}
