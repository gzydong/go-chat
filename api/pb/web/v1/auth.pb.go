// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.26.1
// source: web/v1/auth.proto

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

// 登录接口请求参数
type AuthLoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 登录手机号
	Mobile string `protobuf:"bytes,1,opt,name=mobile,proto3" json:"mobile,omitempty" binding:"required"`
	// 登录密码
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty" binding:"required"`
	// 登录平台
	Platform string `protobuf:"bytes,3,opt,name=platform,proto3" json:"platform,omitempty" binding:"required,oneof=h5 ios windows mac web"`
}

func (x *AuthLoginRequest) Reset() {
	*x = AuthLoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthLoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthLoginRequest) ProtoMessage() {}

func (x *AuthLoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthLoginRequest.ProtoReflect.Descriptor instead.
func (*AuthLoginRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{0}
}

func (x *AuthLoginRequest) GetMobile() string {
	if x != nil {
		return x.Mobile
	}
	return ""
}

func (x *AuthLoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *AuthLoginRequest) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

// 登录接口响应参数
type AuthLoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Token 类型
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// token
	AccessToken string `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	// 过期时间
	ExpiresIn int32 `protobuf:"varint,3,opt,name=expires_in,json=expiresIn,proto3" json:"expires_in,omitempty"`
}

func (x *AuthLoginResponse) Reset() {
	*x = AuthLoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthLoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthLoginResponse) ProtoMessage() {}

func (x *AuthLoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthLoginResponse.ProtoReflect.Descriptor instead.
func (*AuthLoginResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{1}
}

func (x *AuthLoginResponse) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *AuthLoginResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AuthLoginResponse) GetExpiresIn() int32 {
	if x != nil {
		return x.ExpiresIn
	}
	return 0
}

// 注册接口请求参数
type AuthRegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 登录手机号
	Nickname string `protobuf:"bytes,1,opt,name=nickname,proto3" json:"nickname,omitempty" binding:"required,min=2,max=30"`
	// 登录手机号
	Mobile string `protobuf:"bytes,2,opt,name=mobile,proto3" json:"mobile,omitempty" binding:"required,len=11,phone"`
	// 登录密码
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty" binding:"required,min=6,max=16"`
	// 登录平台
	Platform string `protobuf:"bytes,4,opt,name=platform,proto3" json:"platform,omitempty" binding:"required,oneof=h5 ios windows mac web"`
	// 短信验证码
	SmsCode string `protobuf:"bytes,5,opt,name=sms_code,json=smsCode,proto3" json:"sms_code,omitempty" binding:"required"`
}

func (x *AuthRegisterRequest) Reset() {
	*x = AuthRegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRegisterRequest) ProtoMessage() {}

func (x *AuthRegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRegisterRequest.ProtoReflect.Descriptor instead.
func (*AuthRegisterRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{2}
}

func (x *AuthRegisterRequest) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *AuthRegisterRequest) GetMobile() string {
	if x != nil {
		return x.Mobile
	}
	return ""
}

func (x *AuthRegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *AuthRegisterRequest) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *AuthRegisterRequest) GetSmsCode() string {
	if x != nil {
		return x.SmsCode
	}
	return ""
}

// 注册接口响应参数
type AuthRegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AuthRegisterResponse) Reset() {
	*x = AuthRegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRegisterResponse) ProtoMessage() {}

func (x *AuthRegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRegisterResponse.ProtoReflect.Descriptor instead.
func (*AuthRegisterResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{3}
}

// Token 刷新接口请求参数
type AuthRefreshRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AuthRefreshRequest) Reset() {
	*x = AuthRefreshRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRefreshRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRefreshRequest) ProtoMessage() {}

func (x *AuthRefreshRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRefreshRequest.ProtoReflect.Descriptor instead.
func (*AuthRefreshRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{4}
}

// Token 刷新接口响应参数
type AuthRefreshResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	AccessToken string `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	ExpiresIn   int32  `protobuf:"varint,3,opt,name=expires_in,json=expiresIn,proto3" json:"expires_in,omitempty"`
}

func (x *AuthRefreshResponse) Reset() {
	*x = AuthRefreshResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRefreshResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRefreshResponse) ProtoMessage() {}

func (x *AuthRefreshResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRefreshResponse.ProtoReflect.Descriptor instead.
func (*AuthRefreshResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{5}
}

func (x *AuthRefreshResponse) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *AuthRefreshResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AuthRefreshResponse) GetExpiresIn() int32 {
	if x != nil {
		return x.ExpiresIn
	}
	return 0
}

// 找回密码接口请求参数
type AuthForgetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 手机号
	Mobile string `protobuf:"bytes,1,opt,name=mobile,proto3" json:"mobile,omitempty" binding:"required,len=11,phone"`
	// 登录密码
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty" binding:"required,min=6,max=16"`
	// 短信验证码
	SmsCode string `protobuf:"bytes,3,opt,name=sms_code,json=smsCode,proto3" json:"sms_code,omitempty" binding:"required"`
}

func (x *AuthForgetRequest) Reset() {
	*x = AuthForgetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthForgetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthForgetRequest) ProtoMessage() {}

func (x *AuthForgetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthForgetRequest.ProtoReflect.Descriptor instead.
func (*AuthForgetRequest) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{6}
}

func (x *AuthForgetRequest) GetMobile() string {
	if x != nil {
		return x.Mobile
	}
	return ""
}

func (x *AuthForgetRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *AuthForgetRequest) GetSmsCode() string {
	if x != nil {
		return x.SmsCode
	}
	return ""
}

// 找回密码接口响应参数
type AuthForgetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AuthForgetResponse) Reset() {
	*x = AuthForgetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_v1_auth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthForgetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthForgetResponse) ProtoMessage() {}

func (x *AuthForgetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_v1_auth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthForgetResponse.ProtoReflect.Descriptor instead.
func (*AuthForgetResponse) Descriptor() ([]byte, []int) {
	return file_web_v1_auth_proto_rawDescGZIP(), []int{7}
}

var File_web_v1_auth_proto protoreflect.FileDescriptor

var file_web_v1_auth_proto_rawDesc = []byte{
	0x0a, 0x11, 0x77, 0x65, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x77, 0x65, 0x62, 0x1a, 0x13, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72,
	0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xca, 0x01,
	0x0a, 0x10, 0x41, 0x75, 0x74, 0x68, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2f, 0x0a, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67,
	0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x06, 0x6d, 0x6f, 0x62,
	0x69, 0x6c, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x50, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x34, 0x9a, 0x84, 0x9e, 0x03,
	0x2f, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72,
	0x65, 0x64, 0x2c, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x3d, 0x68, 0x35, 0x20, 0x69, 0x6f, 0x73, 0x20,
	0x77, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73, 0x20, 0x6d, 0x61, 0x63, 0x20, 0x77, 0x65, 0x62, 0x22,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x22, 0x69, 0x0a, 0x11, 0x41, 0x75,
	0x74, 0x68, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65,
	0x73, 0x5f, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x73, 0x49, 0x6e, 0x22, 0xdd, 0x02, 0x0a, 0x13, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x40, 0x0a,
	0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x24, 0x9a, 0x84, 0x9e, 0x03, 0x1f, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72,
	0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x6d, 0x69, 0x6e, 0x3d, 0x32, 0x2c, 0x6d, 0x61,
	0x78, 0x3d, 0x33, 0x30, 0x22, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x3c, 0x0a, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x24, 0x9a, 0x84, 0x9e, 0x03, 0x1f, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72,
	0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x6c, 0x65, 0x6e, 0x3d, 0x31, 0x31, 0x2c, 0x70,
	0x68, 0x6f, 0x6e, 0x65, 0x22, 0x52, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x12, 0x40, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x24, 0x9a, 0x84, 0x9e, 0x03, 0x1f, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72,
	0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x6d, 0x69, 0x6e, 0x3d, 0x36, 0x2c, 0x6d, 0x61,
	0x78, 0x3d, 0x31, 0x36, 0x22, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12,
	0x50, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x34, 0x9a, 0x84, 0x9e, 0x03, 0x2f, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a,
	0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x3d,
	0x68, 0x35, 0x20, 0x69, 0x6f, 0x73, 0x20, 0x77, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73, 0x20, 0x6d,
	0x61, 0x63, 0x20, 0x77, 0x65, 0x62, 0x22, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x12, 0x32, 0x0a, 0x08, 0x73, 0x6d, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x52, 0x07, 0x73, 0x6d,
	0x73, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x0a,
	0x12, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x6b, 0x0a, 0x13, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x21,
	0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x49, 0x6e,
	0x22, 0xc7, 0x01, 0x0a, 0x11, 0x41, 0x75, 0x74, 0x68, 0x46, 0x6f, 0x72, 0x67, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3c, 0x0a, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x24, 0x9a, 0x84, 0x9e, 0x03, 0x1f, 0x62, 0x69, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x6c,
	0x65, 0x6e, 0x3d, 0x31, 0x31, 0x2c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x22, 0x52, 0x06, 0x6d, 0x6f,
	0x62, 0x69, 0x6c, 0x65, 0x12, 0x40, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x24, 0x9a, 0x84, 0x9e, 0x03, 0x1f, 0x62, 0x69, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x6d,
	0x69, 0x6e, 0x3d, 0x36, 0x2c, 0x6d, 0x61, 0x78, 0x3d, 0x31, 0x36, 0x22, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x73, 0x6d, 0x73, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0x9a, 0x84, 0x9e, 0x03, 0x12, 0x62,
	0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x22, 0x52, 0x07, 0x73, 0x6d, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x41, 0x75,
	0x74, 0x68, 0x46, 0x6f, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x0c, 0x5a, 0x0a, 0x77, 0x65, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x65, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_web_v1_auth_proto_rawDescOnce sync.Once
	file_web_v1_auth_proto_rawDescData = file_web_v1_auth_proto_rawDesc
)

func file_web_v1_auth_proto_rawDescGZIP() []byte {
	file_web_v1_auth_proto_rawDescOnce.Do(func() {
		file_web_v1_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_web_v1_auth_proto_rawDescData)
	})
	return file_web_v1_auth_proto_rawDescData
}

var file_web_v1_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_web_v1_auth_proto_goTypes = []interface{}{
	(*AuthLoginRequest)(nil),     // 0: web.AuthLoginRequest
	(*AuthLoginResponse)(nil),    // 1: web.AuthLoginResponse
	(*AuthRegisterRequest)(nil),  // 2: web.AuthRegisterRequest
	(*AuthRegisterResponse)(nil), // 3: web.AuthRegisterResponse
	(*AuthRefreshRequest)(nil),   // 4: web.AuthRefreshRequest
	(*AuthRefreshResponse)(nil),  // 5: web.AuthRefreshResponse
	(*AuthForgetRequest)(nil),    // 6: web.AuthForgetRequest
	(*AuthForgetResponse)(nil),   // 7: web.AuthForgetResponse
}
var file_web_v1_auth_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_web_v1_auth_proto_init() }
func file_web_v1_auth_proto_init() {
	if File_web_v1_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_web_v1_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthLoginRequest); i {
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
		file_web_v1_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthLoginResponse); i {
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
		file_web_v1_auth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRegisterRequest); i {
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
		file_web_v1_auth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRegisterResponse); i {
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
		file_web_v1_auth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRefreshRequest); i {
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
		file_web_v1_auth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRefreshResponse); i {
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
		file_web_v1_auth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthForgetRequest); i {
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
		file_web_v1_auth_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthForgetResponse); i {
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
			RawDescriptor: file_web_v1_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_web_v1_auth_proto_goTypes,
		DependencyIndexes: file_web_v1_auth_proto_depIdxs,
		MessageInfos:      file_web_v1_auth_proto_msgTypes,
	}.Build()
	File_web_v1_auth_proto = out.File
	file_web_v1_auth_proto_rawDesc = nil
	file_web_v1_auth_proto_goTypes = nil
	file_web_v1_auth_proto_depIdxs = nil
}
