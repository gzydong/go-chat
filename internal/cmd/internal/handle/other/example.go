package other

import (
	"context"
	"fmt"

	"go-chat/internal/repository/model"
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

	var lastId int

	for {
		var items []*model.TalkRecords

		fmt.Println("LastId :", lastId)
		e.db.Table("talk_records").Where("id > ?", lastId).Order("id asc").Limit(1000).Scan(&items)

		for _, v := range items {
			e.db.Table("talk_records").Where("id = ?", v.Id).UpdateColumn("sequence", e.sequence.Get(ctx, v.UserId, v.ReceiverId))
		}

		if len(items) < 1000 {
			break
		}

		lastId = items[len(items)-1].Id
	}

	fmt.Println("Job ExampleHandle End")

	return nil
}
