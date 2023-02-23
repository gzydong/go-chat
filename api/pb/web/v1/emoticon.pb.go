// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.9
// source: web/v1/emoticon.proto

package web

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

type EmoticonListItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MediaId int32  `protobuf:"varint,1,opt,name=media_id,json=mediaId,proto3" json:"media_id,omitempty"`
	Src     string `protobuf:"bytes,2,opt,name=src,proto3" json:"src,omitempty"`
}

func (x *EmoticonListItem) Reset() {
	*x = EmoticonListItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonListItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonListItem) ProtoMessage() {}

func (x *EmoticonListItem) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonListItem.ProtoReflect.Descriptor instead.
func (*EmoticonListItem) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{0}
}

func (x *EmoticonListItem) GetMediaId() int32 {
	if x != nil {
		return x.MediaId
	}
	return 0
}

func (x *EmoticonListItem) GetSrc() string {
	if x != nil {
		return x.Src
	}
	return ""
}

// 添加或移出表情包接口请求参数
type EmoticonSetSystemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmoticonId int32 `protobuf:"varint,1,opt,name=emoticon_id,json=emoticonId,proto3" json:"emoticon_id,omitempty" binding:"required"`
	Type       int32 `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty" binding:"required,oneof=1 2"`
}

func (x *EmoticonSetSystemRequest) Reset() {
	*x = EmoticonSetSystemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonSetSystemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonSetSystemRequest) ProtoMessage() {}

func (x *EmoticonSetSystemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonSetSystemRequest.ProtoReflect.Descriptor instead.
func (*EmoticonSetSystemRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{1}
}

func (x *EmoticonSetSystemRequest) GetEmoticonId() int32 {
	if x != nil {
		return x.EmoticonId
	}
	return 0
}

func (x *EmoticonSetSystemRequest) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

// 添加或移出表情包接口响应参数
type EmoticonSetSystemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmoticonId int32               `protobuf:"varint,1,opt,name=emoticon_id,json=emoticonId,proto3" json:"emoticon_id,omitempty"`
	Url        string              `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Name       string              `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	List       []*EmoticonListItem `protobuf:"bytes,4,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *EmoticonSetSystemResponse) Reset() {
	*x = EmoticonSetSystemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonSetSystemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonSetSystemResponse) ProtoMessage() {}

func (x *EmoticonSetSystemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonSetSystemResponse.ProtoReflect.Descriptor instead.
func (*EmoticonSetSystemResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{2}
}

func (x *EmoticonSetSystemResponse) GetEmoticonId() int32 {
	if x != nil {
		return x.EmoticonId
	}
	return 0
}

func (x *EmoticonSetSystemResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *EmoticonSetSystemResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EmoticonSetSystemResponse) GetList() []*EmoticonListItem {
	if x != nil {
		return x.List
	}
	return nil
}

// 删除表情包接口请求参数
type EmoticonDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids string `protobuf:"bytes,1,opt,name=ids,proto3" json:"ids,omitempty" form:"ids" binding:"required,ids"`
}

func (x *EmoticonDeleteRequest) Reset() {
	*x = EmoticonDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonDeleteRequest) ProtoMessage() {}

func (x *EmoticonDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonDeleteRequest.ProtoReflect.Descriptor instead.
func (*EmoticonDeleteRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{3}
}

func (x *EmoticonDeleteRequest) GetIds() string {
	if x != nil {
		return x.Ids
	}
	return ""
}

// 删除表情包接口响应参数
type EmoticonDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmoticonDeleteResponse) Reset() {
	*x = EmoticonDeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonDeleteResponse) ProtoMessage() {}

func (x *EmoticonDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonDeleteResponse.ProtoReflect.Descriptor instead.
func (*EmoticonDeleteResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{4}
}

// 系统表情包列表接口请求参数
type EmoticonSysListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmoticonSysListRequest) Reset() {
	*x = EmoticonSysListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonSysListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonSysListRequest) ProtoMessage() {}

func (x *EmoticonSysListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonSysListRequest.ProtoReflect.Descriptor instead.
func (*EmoticonSysListRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{5}
}

// 系统表情包列表接口响应参数
type EmoticonSysListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*EmoticonSysListResponse_Item `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *EmoticonSysListResponse) Reset() {
	*x = EmoticonSysListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonSysListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonSysListResponse) ProtoMessage() {}

func (x *EmoticonSysListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonSysListResponse.ProtoReflect.Descriptor instead.
func (*EmoticonSysListResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{6}
}

func (x *EmoticonSysListResponse) GetItems() []*EmoticonSysListResponse_Item {
	if x != nil {
		return x.Items
	}
	return nil
}

// 用户表情包列表接口请求参数
type EmoticonListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmoticonListRequest) Reset() {
	*x = EmoticonListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonListRequest) ProtoMessage() {}

func (x *EmoticonListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonListRequest.ProtoReflect.Descriptor instead.
func (*EmoticonListRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{7}
}

// 用户表情包列表接口响应参数
type EmoticonListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SysEmoticon     []*EmoticonListResponse_SysEmoticon `protobuf:"bytes,1,rep,name=sys_emoticon,json=sysEmoticon,proto3" json:"sys_emoticon,omitempty"`
	CollectEmoticon []*EmoticonListItem                 `protobuf:"bytes,2,rep,name=collect_emoticon,json=collectEmoticon,proto3" json:"collect_emoticon,omitempty"`
}

func (x *EmoticonListResponse) Reset() {
	*x = EmoticonListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonListResponse) ProtoMessage() {}

func (x *EmoticonListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonListResponse.ProtoReflect.Descriptor instead.
func (*EmoticonListResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{8}
}

func (x *EmoticonListResponse) GetSysEmoticon() []*EmoticonListResponse_SysEmoticon {
	if x != nil {
		return x.SysEmoticon
	}
	return nil
}

func (x *EmoticonListResponse) GetCollectEmoticon() []*EmoticonListItem {
	if x != nil {
		return x.CollectEmoticon
	}
	return nil
}

// 表情包上传接口请求参数
type EmoticonUploadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmoticonUploadRequest) Reset() {
	*x = EmoticonUploadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonUploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonUploadRequest) ProtoMessage() {}

func (x *EmoticonUploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonUploadRequest.ProtoReflect.Descriptor instead.
func (*EmoticonUploadRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{9}
}

// 表情包上传接口响应参数
type EmoticonUploadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MediaId int32  `protobuf:"varint,1,opt,name=media_id,json=mediaId,proto3" json:"media_id,omitempty"`
	Src     string `protobuf:"bytes,2,opt,name=src,proto3" json:"src,omitempty"`
}

func (x *EmoticonUploadResponse) Reset() {
	*x = EmoticonUploadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonUploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonUploadResponse) ProtoMessage() {}

func (x *EmoticonUploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonUploadResponse.ProtoReflect.Descriptor instead.
func (*EmoticonUploadResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{10}
}

func (x *EmoticonUploadResponse) GetMediaId() int32 {
	if x != nil {
		return x.MediaId
	}
	return 0
}

func (x *EmoticonUploadResponse) GetSrc() string {
	if x != nil {
		return x.Src
	}
	return ""
}

type EmoticonSysListResponse_Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Icon   string `protobuf:"bytes,3,opt,name=icon,proto3" json:"icon,omitempty"`
	Status int32  `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *EmoticonSysListResponse_Item) Reset() {
	*x = EmoticonSysListResponse_Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonSysListResponse_Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonSysListResponse_Item) ProtoMessage() {}

func (x *EmoticonSysListResponse_Item) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonSysListResponse_Item.ProtoReflect.Descriptor instead.
func (*EmoticonSysListResponse_Item) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{6, 0}
}

func (x *EmoticonSysListResponse_Item) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *EmoticonSysListResponse_Item) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EmoticonSysListResponse_Item) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

func (x *EmoticonSysListResponse_Item) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type EmoticonListResponse_SysEmoticon struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmoticonId int32               `protobuf:"varint,1,opt,name=emoticon_id,json=emoticonId,proto3" json:"emoticon_id,omitempty"`
	Url        string              `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Name       string              `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	List       []*EmoticonListItem `protobuf:"bytes,4,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *EmoticonListResponse_SysEmoticon) Reset() {
	*x = EmoticonListResponse_SysEmoticon{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_emoticon_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoticonListResponse_SysEmoticon) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoticonListResponse_SysEmoticon) ProtoMessage() {}

func (x *EmoticonListResponse_SysEmoticon) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_emoticon_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoticonListResponse_SysEmoticon.ProtoReflect.Descriptor instead.
func (*EmoticonListResponse_SysEmoticon) Descriptor() ([]byte, []int) {
	return file_web_v1_emoticon_proto_rawDescGZIP(), []int{8, 0}
}

func (x *EmoticonListResponse_SysEmoticon) GetEmoticonId() int32 {
	if x != nil {
		return x.EmoticonId
	}
	return 0
}

func (x *EmoticonListResponse_SysEmoticon) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *EmoticonListResponse_SysEmoticon) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EmoticonListResponse_SysEmoticon) GetList() []*EmoticonListItem {
	if x != nil {
		return x.List
	}
	return nil
}

var File_web_v1_emoticon_proto protoreflect.FileDescriptor

var file_web_v1_emoticon_proto_rawDesc = []byte{
	0x0a, 0x15, 0x77, 0x65, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x77, 0x65, 0x62, 0x1a, 0x13, 0x74, 0x61,
	0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3f, 0x0a, 0x10, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x4c, 0x69, 0x73,
	0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x73, 0x72, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73,
	0x72, 0x63, 0x22, 0x8b, 0x01, 0x0a, 0x18, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x53,
	0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x38, 0x0a, 0x0b, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x0a, 0x65,
	0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x21, 0x9a, 0x84, 0x9e, 0x03, 0x1c, 0x62, 0x69,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c,
	0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x3d, 0x31, 0x20, 0x32, 0x22, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x22, 0x8d, 0x01, 0x0a, 0x19, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x53, 0x65, 0x74,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x2e, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63,
	0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74,
	0x22, 0x51, 0x0a, 0x15, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x03, 0x69, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x26, 0x9a, 0x84, 0x9e, 0x03, 0x21, 0x66, 0x6f, 0x72,
	0x6d, 0x3a, 0x22, 0x69, 0x64, 0x73, 0x22, 0x20, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a,
	0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x69, 0x64, 0x73, 0x22, 0x52, 0x03,
	0x69, 0x64, 0x73, 0x22, 0x18, 0x0a, 0x16, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x0a,
	0x16, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x53, 0x79, 0x73, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xaa, 0x01, 0x0a, 0x17, 0x45, 0x6d, 0x6f, 0x74,
	0x69, 0x63, 0x6f, 0x6e, 0x53, 0x79, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x21, 0x2e, 0x77, 0x65, 0x62, 0x2e, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f,
	0x6e, 0x53, 0x79, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x1a, 0x56, 0x0a, 0x04,
	0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x15, 0x0a, 0x13, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xa3, 0x02, 0x0a, 0x14,
	0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x0c, 0x73, 0x79, 0x73, 0x5f, 0x65, 0x6d, 0x6f, 0x74,
	0x69, 0x63, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x77, 0x65, 0x62,
	0x2e, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x79, 0x73, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f,
	0x6e, 0x52, 0x0b, 0x73, 0x79, 0x73, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x12, 0x40,
	0x0a, 0x10, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x5f, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x63,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x2e, 0x45,
	0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x0f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e,
	0x1a, 0x7f, 0x0a, 0x0b, 0x53, 0x79, 0x73, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x12,
	0x1f, 0x0a, 0x0b, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x2e, 0x45, 0x6d, 0x6f, 0x74, 0x69,
	0x63, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x6c, 0x69, 0x73,
	0x74, 0x22, 0x17, 0x0a, 0x15, 0x45, 0x6d, 0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x45, 0x0a, 0x16, 0x45, 0x6d,
	0x6f, 0x74, 0x69, 0x63, 0x6f, 0x6e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x72, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x72,
	0x63, 0x42, 0x0c, 0x5a, 0x0a, 0x77, 0x65, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x65, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_web_v1_emoticon_proto_rawDescOnce sync.Once
	file_web_v1_emoticon_proto_rawDescData = file_web_v1_emoticon_proto_rawDesc
)

func file_web_v1_emoticon_proto_rawDescGZIP() []byte {
	file_web_v1_emoticon_proto_rawDescOnce.Do(func() {
		file_web_v1_emoticon_proto_rawDescData = protoimpl.X.CompressGZIP(file_web_v1_emoticon_proto_rawDescData)
	})
	return file_web_v1_emoticon_proto_rawDescData
}

var file_web_v1_emoticon_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_web_v1_emoticon_proto_goTypes = []any{
	(*EmoticonListItem)(nil),                 // 0: web.EmoticonListItem
	(*EmoticonSetSystemRequest)(nil),         // 1: web.EmoticonSetSystemRequest
	(*EmoticonSetSystemResponse)(nil),        // 2: web.EmoticonSetSystemResponse
	(*EmoticonDeleteRequest)(nil),            // 3: web.EmoticonDeleteRequest
	(*EmoticonDeleteResponse)(nil),           // 4: web.EmoticonDeleteResponse
	(*EmoticonSysListRequest)(nil),           // 5: web.EmoticonSysListRequest
	(*EmoticonSysListResponse)(nil),          // 6: web.EmoticonSysListResponse
	(*EmoticonListRequest)(nil),              // 7: web.EmoticonListRequest
	(*EmoticonListResponse)(nil),             // 8: web.EmoticonListResponse
	(*EmoticonUploadRequest)(nil),            // 9: web.EmoticonUploadRequest
	(*EmoticonUploadResponse)(nil),           // 10: web.EmoticonUploadResponse
	(*EmoticonSysListResponse_Item)(nil),     // 11: web.EmoticonSysListResponse.Item
	(*EmoticonListResponse_SysEmoticon)(nil), // 12: web.EmoticonListResponse.SysEmoticon
}
var file_web_v1_emoticon_proto_depIdxs = []int32{
	0,  // 0: web.EmoticonSetSystemResponse.list:type_name -> web.EmoticonListItem
	11, // 1: web.EmoticonSysListResponse.items:type_name -> web.EmoticonSysListResponse.Item
	12, // 2: web.EmoticonListResponse.sys_emoticon:type_name -> web.EmoticonListResponse.SysEmoticon
	0,  // 3: web.EmoticonListResponse.collect_emoticon:type_name -> web.EmoticonListItem
	0,  // 4: web.EmoticonListResponse.SysEmoticon.list:type_name -> web.EmoticonListItem
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_web_v1_emoticon_proto_init() }
func file_web_v1_emoticon_proto_init() {
	if File_web_v1_emoticon_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_web_v1_emoticon_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonListItem); i {
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
		file_web_v1_emoticon_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonSetSystemRequest); i {
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
		file_web_v1_emoticon_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonSetSystemResponse); i {
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
		file_web_v1_emoticon_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonDeleteRequest); i {
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
		file_web_v1_emoticon_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonDeleteResponse); i {
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
		file_web_v1_emoticon_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonSysListRequest); i {
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
		file_web_v1_emoticon_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonSysListResponse); i {
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
		file_web_v1_emoticon_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonListRequest); i {
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
		file_web_v1_emoticon_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonListResponse); i {
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
		file_web_v1_emoticon_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonUploadRequest); i {
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
		file_web_v1_emoticon_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonUploadResponse); i {
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
		file_web_v1_emoticon_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonSysListResponse_Item); i {
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
		file_web_v1_emoticon_proto_msgTypes[12].Exporter = func(v any, i int) any {
			switch v := v.(*EmoticonListResponse_SysEmoticon); i {
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
			RawDescriptor: file_web_v1_emoticon_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_web_v1_emoticon_proto_goTypes,
		DependencyIndexes: file_web_v1_emoticon_proto_depIdxs,
		MessageInfos:      file_web_v1_emoticon_proto_msgTypes,
	}.Build()
	File_web_v1_emoticon_proto = out.File
	file_web_v1_emoticon_proto_rawDesc = nil
	file_web_v1_emoticon_proto_goTypes = nil
	file_web_v1_emoticon_proto_depIdxs = nil
}
