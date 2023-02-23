package v1

import (
	"go-chat/api/pb/open/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
)

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (c *Index) Index(ctx *ichat.Context) error {

	params := &open.IndexRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	return ctx.Success(params)
}

func (c *Index) Proto(ctx *ichat.Context) error {
	return ctx.Success(&web.AuthLoginResponse{
		Type:        "Type",
		AccessToken: "AccessToken",
		ExpiresIn:   0,
	})
}
