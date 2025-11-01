package v1

import (
	"github.com/gzydong/go-chat/api/pb/open/v1"
	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/pkg/core"
)

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (c *Index) Index(ctx *core.Context) error {

	var in open.IndexRequest
	if err := ctx.Context.ShouldBind(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	return ctx.Success(&in)
}

func (c *Index) Proto(ctx *core.Context) error {
	return ctx.Success(&web.AuthLoginResponse{
		Type:        "Type",
		AccessToken: "AccessToken",
		ExpiresIn:   0,
	})
}
