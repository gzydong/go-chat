package entity

import (
	"go-chat/internal/pkg/core/errorx"
)

var (
	ErrDataNotFound           = errorx.New(404, "访问数据不存在")
	ErrPermissionDenied       = errorx.New(403, "无权访问资源")
	ErrTooFrequentOperation   = errorx.New(429, "操作频繁请稍后再试")
	ErrInvalidParams          = errorx.New(400, "参数错误")
	ErrUserNotExist           = errorx.New(100004, "用户不存在")
	ErrAccountOrPassword      = errorx.New(100005, "账号不存在或密码错误")
	ErrPhoneExist             = errorx.New(100006, "手机号已存在")
	ErrAccountOrPasswordError = errorx.New(100007, "账号密码填写错误")
	ErrSmsCodeError           = errorx.New(100008, "短信验证码填写错误")
	ErrAccountDisabled        = errorx.New(100009, "账号已被管理员禁用，如有问题请联系管理员！")
	ErrGroupDismissed         = errorx.New(110001, "群组已解散")
	ErrGroupMemberLimit       = errorx.New(110002, "群成员数量已达到上限")
	ErrGroupNotExist          = errorx.New(110003, "群组不存在")
)
