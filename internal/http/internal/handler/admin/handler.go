package admin

import "go-chat/internal/http/internal/handler/admin/v1"

type V1 struct {
	Index *v1.Index
	Auth  *v1.Auth
}

type V2 struct {
}

type Handler struct {
	V1 *V1
	V2 *V2
}
