package entity

import "errors"

var (
	// ErrPermissionDenied 无权访问资源
	ErrPermissionDenied = errors.New("无权限访问！")
)
