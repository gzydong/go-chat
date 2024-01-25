package v1

import "go-chat/internal/pkg/core"

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (c *Index) Index(ctx *core.Context) error {
	return nil
}
