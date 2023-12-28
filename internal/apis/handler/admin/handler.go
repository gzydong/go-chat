package admin

import (
	v12 "go-chat/internal/apis/handler/admin/v1"
)

type V1 struct {
	Index *v12.Index
	Auth  *v12.Auth
}

type V2 struct{}

type Handler struct {
	V1 *V1
	V2 *V2
}
