package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecordsDelete struct {
	ichat.Repo[model.TalkRecordsDelete]
}

func NewTalkRecordsDelete(db *gorm.DB) *TalkRecordsDelete {
	return &TalkRecordsDelete{Repo: ichat.NewRepo[model.TalkRecordsDelete](db)}
}

func (t *TalkRecordsDelete) FindAllRecordIds(ctx context.Context, ids []int, userId int) ([]int, error) {
	var records []int

	err := t.Db.WithContext(ctx).Table("talk_records_delete").Select("record_id").Where("user_id =?", userId).Where("record_id IN (?)", ids).Scan(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
