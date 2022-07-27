package v1

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
)

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (c *Index) Index(ctx *ichat.Context) error {
	return ctx.Success(entity.H{
		"uid":      ctx.UserId(),
		"is_guest": ctx.IsGuest(),
	})
}

func (c *Index) Proto(ctx *ichat.Context) error {
	return ctx.Success(web.AuthLoginResponse{
		Type:        "Type",
		AccessToken: "AccessToken",
		ExpiresIn:   0,
	})
}
