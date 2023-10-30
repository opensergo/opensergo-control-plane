// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: rate_limit.proto

package __

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/anypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RateLimitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Threshold int64 `protobuf:"varint,1,opt,name=threshold,proto3" json:"threshold,omitempty"`
}

func (x *RateLimitResponse) Reset() {
	*x = RateLimitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rate_limit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitResponse) ProtoMessage() {}

func (x *RateLimitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rate_limit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitResponse.ProtoReflect.Descriptor instead.
func (*RateLimitResponse) Descriptor() ([]byte, []int) {
	return file_rate_limit_proto_rawDescGZIP(), []int{0}
}

func (x *RateLimitResponse) GetThreshold() int64 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

type RateLimitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Threshold int64 `protobuf:"varint,1,opt,name=threshold,proto3" json:"threshold,omitempty"`
}

func (x *RateLimitRequest) Reset() {
	*x = RateLimitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rate_limit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitRequest) ProtoMessage() {}

func (x *RateLimitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rate_limit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitRequest.ProtoReflect.Descriptor instead.
func (*RateLimitRequest) Descriptor() ([]byte, []int) {
	return file_rate_limit_proto_rawDescGZIP(), []int{1}
}

func (x *RateLimitRequest) GetThreshold() int64 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

var File_rate_limit_proto protoreflect.FileDescriptor

var file_rate_limit_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x24, 0x69, 0x6f, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x72, 0x67, 0x6f,
	0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x61,
	0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x31, 0x0a, 0x11, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65,
	0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x68, 0x72,
	0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x22, 0x30, 0x0a, 0x10, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x68,
	0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74,
	0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x32, 0x90, 0x01, 0x0a, 0x10, 0x52, 0x61, 0x74,
	0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7c, 0x0a,
	0x09, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x36, 0x2e, 0x69, 0x6f, 0x2e,
	0x6f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x72, 0x67, 0x6f, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x37, 0x2e, 0x69, 0x6f, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x72, 0x67,
	0x6f, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x72,
	0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x04, 0x5a, 0x02, 0x2e,
	0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rate_limit_proto_rawDescOnce sync.Once
	file_rate_limit_proto_rawDescData = file_rate_limit_proto_rawDesc
)

func file_rate_limit_proto_rawDescGZIP() []byte {
	file_rate_limit_proto_rawDescOnce.Do(func() {
		file_rate_limit_proto_rawDescData = protoimpl.X.CompressGZIP(file_rate_limit_proto_rawDescData)
	})
	return file_rate_limit_proto_rawDescData
}

var file_rate_limit_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rate_limit_proto_goTypes = []interface{}{
	(*RateLimitResponse)(nil), // 0: io.opensergo.plugin.proto.rate_limit.RateLimitResponse
	(*RateLimitRequest)(nil),  // 1: io.opensergo.plugin.proto.rate_limit.RateLimitRequest
}
var file_rate_limit_proto_depIdxs = []int32{
	1, // 0: io.opensergo.plugin.proto.rate_limit.RateLimitService.RateLimit:input_type -> io.opensergo.plugin.proto.rate_limit.RateLimitRequest
	0, // 1: io.opensergo.plugin.proto.rate_limit.RateLimitService.RateLimit:output_type -> io.opensergo.plugin.proto.rate_limit.RateLimitResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rate_limit_proto_init() }
func file_rate_limit_proto_init() {
	if File_rate_limit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rate_limit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimitResponse); i {
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
		file_rate_limit_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimitRequest); i {
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
			RawDescriptor: file_rate_limit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rate_limit_proto_goTypes,
		DependencyIndexes: file_rate_limit_proto_depIdxs,
		MessageInfos:      file_rate_limit_proto_msgTypes,
	}.Build()
	File_rate_limit_proto = out.File
	file_rate_limit_proto_rawDesc = nil
	file_rate_limit_proto_goTypes = nil
	file_rate_limit_proto_depIdxs = nil
}