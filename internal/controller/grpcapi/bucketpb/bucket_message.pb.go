// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: bucket_message.proto

package bucketpb

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Ip       string `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bucket_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_bucket_message_proto_msgTypes[0]
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
	return file_bucket_message_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *Request) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Request) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

var File_bucket_message_proto protoreflect.FileDescriptor

var file_bucket_message_proto_rawDesc = []byte{
	0x0a, 0x14, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x4b,
	0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67,
	0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x42, 0x26, 0x5a, 0x24, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bucket_message_proto_rawDescOnce sync.Once
	file_bucket_message_proto_rawDescData = file_bucket_message_proto_rawDesc
)

func file_bucket_message_proto_rawDescGZIP() []byte {
	file_bucket_message_proto_rawDescOnce.Do(func() {
		file_bucket_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_bucket_message_proto_rawDescData)
	})
	return file_bucket_message_proto_rawDescData
}

var file_bucket_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_bucket_message_proto_goTypes = []interface{}{
	(*Request)(nil), // 0: bucket.Request
}
var file_bucket_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_bucket_message_proto_init() }
func file_bucket_message_proto_init() {
	if File_bucket_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bucket_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_bucket_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bucket_message_proto_goTypes,
		DependencyIndexes: file_bucket_message_proto_depIdxs,
		MessageInfos:      file_bucket_message_proto_msgTypes,
	}.Build()
	File_bucket_message_proto = out.File
	file_bucket_message_proto_rawDesc = nil
	file_bucket_message_proto_goTypes = nil
	file_bucket_message_proto_depIdxs = nil
}
