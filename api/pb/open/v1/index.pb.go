// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: open/v1/index.proto

package open

import (
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
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

// 分配管理员接口请求参数
type IndexRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId int32 `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty" binding:"required"`
	UserId  int32 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" binding:"required"`
	Mode    int32 `protobuf:"varint,3,opt,name=mode,proto3" json:"mode,omitempty" binding:"required,oneof=1 2"`
}

func (x *IndexRequest) Reset() {
	*x = IndexRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_open_v1_index_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexRequest) ProtoMessage() {}

func (x *IndexRequest) ProtoReflect() protoreflect.Message {
	mi := &file_open_v1_index_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexRequest.ProtoReflect.Descriptor instead.
func (*IndexRequest) Descriptor() ([]byte, []int) {
	return file_open_v1_index_proto_rawDescGZIP(), []int{0}
}

func (x *IndexRequest) GetGroupId() int32 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *IndexRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *IndexRequest) GetMode() int32 {
	if x != nil {
		return x.Mode
	}
	return 0
}

var File_open_v1_index_proto protoreflect.FileDescriptor

var file_open_v1_index_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6f, 0x70, 0x65, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x77, 0x65, 0x62, 0x1a, 0x13, 0x74, 0x61, 0x67, 0x67,
	0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xab, 0x01, 0x0a, 0x0c, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x32, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67,
	0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x07, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x21, 0x9a, 0x84, 0x9e, 0x03, 0x1c, 0x62, 0x69, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x6f, 0x6e, 0x65,
	0x6f, 0x66, 0x3d, 0x31, 0x20, 0x32, 0x22, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x42, 0x0e, 0x5a,
	0x0c, 0x6f, 0x70, 0x65, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x6f, 0x70, 0x65, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_open_v1_index_proto_rawDescOnce sync.Once
	file_open_v1_index_proto_rawDescData = file_open_v1_index_proto_rawDesc
)

func file_open_v1_index_proto_rawDescGZIP() []byte {
	file_open_v1_index_proto_rawDescOnce.Do(func() {
		file_open_v1_index_proto_rawDescData = protoimpl.X.CompressGZIP(file_open_v1_index_proto_rawDescData)
	})
	return file_open_v1_index_proto_rawDescData
}

var file_open_v1_index_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_open_v1_index_proto_goTypes = []interface{}{
	(*IndexRequest)(nil), // 0: web.IndexRequest
}
var file_open_v1_index_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_open_v1_index_proto_init() }
func file_open_v1_index_proto_init() {
	if File_open_v1_index_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_open_v1_index_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexRequest); i {
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
			RawDescriptor: file_open_v1_index_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_open_v1_index_proto_goTypes,
		DependencyIndexes: file_open_v1_index_proto_depIdxs,
		MessageInfos:      file_open_v1_index_proto_msgTypes,
	}.Build()
	File_open_v1_index_proto = out.File
	file_open_v1_index_proto_rawDesc = nil
	file_open_v1_index_proto_goTypes = nil
	file_open_v1_index_proto_depIdxs = nil
}
