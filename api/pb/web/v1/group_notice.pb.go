// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: web/v1/group_notice.proto

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

type GroupNoticeDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId  int32 `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty" binding:"required"`
	NoticeId int32 `protobuf:"varint,2,opt,name=notice_id,json=noticeId,proto3" json:"notice_id,omitempty" binding:"required"`
}

func (x *GroupNoticeDeleteRequest) Reset() {
	*x = GroupNoticeDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_group_notice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupNoticeDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupNoticeDeleteRequest) ProtoMessage() {}

func (x *GroupNoticeDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_group_notice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupNoticeDeleteRequest.ProtoReflect.Descriptor instead.
func (*GroupNoticeDeleteRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_group_notice_proto_rawDescGZIP(), []int{0}
}

func (x *GroupNoticeDeleteRequest) GetGroupId() int32 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *GroupNoticeDeleteRequest) GetNoticeId() int32 {
	if x != nil {
		return x.NoticeId
	}
	return 0
}

type GroupNoticeDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GroupNoticeDeleteResponse) Reset() {
	*x = GroupNoticeDeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_group_notice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupNoticeDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupNoticeDeleteResponse) ProtoMessage() {}

func (x *GroupNoticeDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_group_notice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupNoticeDeleteResponse.ProtoReflect.Descriptor instead.
func (*GroupNoticeDeleteResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_group_notice_proto_rawDescGZIP(), []int{1}
}

type GroupNoticeEditRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId int32  `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty" binding:"required"`
	Content string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty" binding:"required"`
}

func (x *GroupNoticeEditRequest) Reset() {
	*x = GroupNoticeEditRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_group_notice_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupNoticeEditRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupNoticeEditRequest) ProtoMessage() {}

func (x *GroupNoticeEditRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_group_notice_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupNoticeEditRequest.ProtoReflect.Descriptor instead.
func (*GroupNoticeEditRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_group_notice_proto_rawDescGZIP(), []int{2}
}

func (x *GroupNoticeEditRequest) GetGroupId() int32 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *GroupNoticeEditRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type GroupNoticeEditResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GroupNoticeEditResponse) Reset() {
	*x = GroupNoticeEditResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_group_notice_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupNoticeEditResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupNoticeEditResponse) ProtoMessage() {}

func (x *GroupNoticeEditResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_group_notice_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupNoticeEditResponse.ProtoReflect.Descriptor instead.
func (*GroupNoticeEditResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_group_notice_proto_rawDescGZIP(), []int{3}
}

type GroupNoticeListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId int32 `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty" form:"group_id" binding:"required"`
}

func (x *GroupNoticeListRequest) Reset() {
	*x = GroupNoticeListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_group_notice_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupNoticeListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupNoticeListRequest) ProtoMessage() {}

func (x *GroupNoticeListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_group_notice_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupNoticeListRequest.ProtoReflect.Descriptor instead.
func (*GroupNoticeListRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_group_notice_proto_rawDescGZIP(), []int{4}
}

func (x *GroupNoticeListRequest) GetGroupId() int32 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

type GroupNoticeListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*GroupNoticeListResponse_Item `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GroupNoticeListResponse) Reset() {
	*x = GroupNoticeListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_group_notice_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupNoticeListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupNoticeListResponse) ProtoMessage() {}

func (x *GroupNoticeListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_group_notice_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupNoticeListResponse.ProtoReflect.Descriptor instead.
func (*GroupNoticeListResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_group_notice_proto_rawDescGZIP(), []int{5}
}

func (x *GroupNoticeListResponse) GetItems() []*GroupNoticeListResponse_Item {
	if x != nil {
		return x.Items
	}
	return nil
}

type GroupNoticeListResponse_Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title        string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Content      string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	IsTop        int32  `protobuf:"varint,4,opt,name=is_top,json=isTop,proto3" json:"is_top,omitempty"`
	IsConfirm    int32  `protobuf:"varint,5,opt,name=is_confirm,json=isConfirm,proto3" json:"is_confirm,omitempty"`
	ConfirmUsers string `protobuf:"bytes,6,opt,name=confirm_users,json=confirmUsers,proto3" json:"confirm_users,omitempty"`
	Avatar       string `protobuf:"bytes,7,opt,name=avatar,proto3" json:"avatar,omitempty"`
	CreatorId    int32  `protobuf:"varint,8,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	CreatedAt    string `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt    string `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *GroupNoticeListResponse_Item) Reset() {
	*x = GroupNoticeListResponse_Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_group_notice_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupNoticeListResponse_Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupNoticeListResponse_Item) ProtoMessage() {}

func (x *GroupNoticeListResponse_Item) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_group_notice_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupNoticeListResponse_Item.ProtoReflect.Descriptor instead.
func (*GroupNoticeListResponse_Item) Descriptor() ([]byte, []int) {
	return file_web_v1_group_notice_proto_rawDescGZIP(), []int{5, 0}
}

func (x *GroupNoticeListResponse_Item) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GroupNoticeListResponse_Item) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *GroupNoticeListResponse_Item) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *GroupNoticeListResponse_Item) GetIsTop() int32 {
	if x != nil {
		return x.IsTop
	}
	return 0
}

func (x *GroupNoticeListResponse_Item) GetIsConfirm() int32 {
	if x != nil {
		return x.IsConfirm
	}
	return 0
}

func (x *GroupNoticeListResponse_Item) GetConfirmUsers() string {
	if x != nil {
		return x.ConfirmUsers
	}
	return ""
}

func (x *GroupNoticeListResponse_Item) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *GroupNoticeListResponse_Item) GetCreatorId() int32 {
	if x != nil {
		return x.CreatorId
	}
	return 0
}

func (x *GroupNoticeListResponse_Item) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *GroupNoticeListResponse_Item) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

var File_web_v1_group_notice_proto protoreflect.FileDescriptor

var file_web_v1_group_notice_proto_rawDesc = []byte{
	0x0a, 0x19, 0x77, 0x65, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x77, 0x65, 0x62,
	0x1a, 0x13, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x18, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x32, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x07, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x09, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12,
	0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x22, 0x52, 0x08, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x49, 0x64, 0x22, 0x1b, 0x0a, 0x19,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x7f, 0x0a, 0x16, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x45, 0x64, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x07,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62,
	0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x22, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x19, 0x0a, 0x17, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x45, 0x64, 0x69, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5c, 0x0a, 0x16, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x6f,
	0x74, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x42, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x42, 0x27, 0x9a, 0x84, 0x9e, 0x03, 0x22, 0x66, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x22, 0x20, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a,
	0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x49, 0x64, 0x22, 0xeb, 0x02, 0x0a, 0x17, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x6f, 0x74,
	0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x37, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x77, 0x65, 0x62, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x1a, 0x96, 0x02, 0x0a, 0x04, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x15, 0x0a, 0x06, 0x69, 0x73, 0x5f, 0x74, 0x6f, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x69, 0x73, 0x54, 0x6f, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x69, 0x73,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x72, 0x6d, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76,
	0x61, 0x74, 0x61, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x42, 0x0c, 0x5a, 0x0a, 0x77, 0x65, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x65, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_web_v1_group_notice_proto_rawDescOnce sync.Once
	file_web_v1_group_notice_proto_rawDescData = file_web_v1_group_notice_proto_rawDesc
)

func file_web_v1_group_notice_proto_rawDescGZIP() []byte {
	file_web_v1_group_notice_proto_rawDescOnce.Do(func() {
		file_web_v1_group_notice_proto_rawDescData = protoimpl.X.CompressGZIP(file_web_v1_group_notice_proto_rawDescData)
	})
	return file_web_v1_group_notice_proto_rawDescData
}

var file_web_v1_group_notice_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_web_v1_group_notice_proto_goTypes = []any{
	(*GroupNoticeDeleteRequest)(nil),     // 0: web.GroupNoticeDeleteRequest
	(*GroupNoticeDeleteResponse)(nil),    // 1: web.GroupNoticeDeleteResponse
	(*GroupNoticeEditRequest)(nil),       // 2: web.GroupNoticeEditRequest
	(*GroupNoticeEditResponse)(nil),      // 3: web.GroupNoticeEditResponse
	(*GroupNoticeListRequest)(nil),       // 4: web.GroupNoticeListRequest
	(*GroupNoticeListResponse)(nil),      // 5: web.GroupNoticeListResponse
	(*GroupNoticeListResponse_Item)(nil), // 6: web.GroupNoticeListResponse.Item
}
var file_web_v1_group_notice_proto_depIdxs = []int32{
	6, // 0: web.GroupNoticeListResponse.items:type_name -> web.GroupNoticeListResponse.Item
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_web_v1_group_notice_proto_init() }
func file_web_v1_group_notice_proto_init() {
	if File_web_v1_group_notice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_web_v1_group_notice_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GroupNoticeDeleteRequest); i {
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
		file_web_v1_group_notice_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GroupNoticeDeleteResponse); i {
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
		file_web_v1_group_notice_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GroupNoticeEditRequest); i {
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
		file_web_v1_group_notice_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GroupNoticeEditResponse); i {
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
		file_web_v1_group_notice_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*GroupNoticeListRequest); i {
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
		file_web_v1_group_notice_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*GroupNoticeListResponse); i {
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
		file_web_v1_group_notice_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GroupNoticeListResponse_Item); i {
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
			RawDescriptor: file_web_v1_group_notice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_web_v1_group_notice_proto_goTypes,
		DependencyIndexes: file_web_v1_group_notice_proto_depIdxs,
		MessageInfos:      file_web_v1_group_notice_proto_msgTypes,
	}.Build()
	File_web_v1_group_notice_proto = out.File
	file_web_v1_group_notice_proto_rawDesc = nil
	file_web_v1_group_notice_proto_goTypes = nil
	file_web_v1_group_notice_proto_depIdxs = nil
}
