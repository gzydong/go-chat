package other

import (
	"context"
	"fmt"

	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type ExampleHandle struct {
	db       *gorm.DB
	sequence *repo.Sequence
}

func NewExampleHandle(db *gorm.DB, sequence *repo.Sequence) *ExampleHandle {
	return &ExampleHandle{db, sequence}
}

func (e *ExampleHandle) Handle(ctx context.Context) error {

	fmt.Println("Job ExampleHandle Start")

	fmt.Println("Job ExampleHandle End")

	return nil
}
