package open

import (
	"go-chat/internal/apis/handler/open/v1"
)

type V1 struct {
	Index *v1.Index
}

type Handler struct {
	V1 *V1
}
