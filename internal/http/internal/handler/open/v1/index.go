package v1

import (
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
		"uid":     ctx.LoginUID(),
		"IsGuest": ctx.IsGuest(),
	})
}
