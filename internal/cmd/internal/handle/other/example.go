package other

import (
	"context"

	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ExampleHandle struct {
	db *gorm.DB
}

func NewExampleHandle(db *gorm.DB) *ExampleHandle {
	return &ExampleHandle{db}
}

func (e *ExampleHandle) Handle(ctx context.Context) error {

	id := 0
	for {
		items := make([]*model.TalkRecords, 0)
		e.db.Table("talk_records").Where("id > ?", id).Where("msg_id = ''").Order("id asc").Limit(100).Scan(&items)

		for _, item := range items {
			e.db.Table("talk_records").Where("id = ?", item.Id).UpdateColumn("msg_id", strutil.NewUuid())
		}

		if len(items) < 100 {
			break
		}

		id = items[len(items)-1].Id
	}

	return nil
}
