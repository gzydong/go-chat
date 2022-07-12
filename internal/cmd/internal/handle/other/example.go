package other

import (
	"context"
)

type ExampleHandle struct {
}

func NewExampleHandle() *ExampleHandle {
	return &ExampleHandle{}
}

func (e *ExampleHandle) Handle(ctx context.Context) error {

	return nil
}
