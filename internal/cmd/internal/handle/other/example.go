package other

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ExampleHandle struct {
	db *gorm.DB
}

func NewExampleHandle(db *gorm.DB) *ExampleHandle {
	return &ExampleHandle{db}
}

func (e *ExampleHandle) Handle(ctx context.Context) error {

	fmt.Println("Job ExampleHandle")

	return nil
}
