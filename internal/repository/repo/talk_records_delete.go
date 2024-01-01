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

func (t *TalkRecordsDelete) FindAllMsgIds(ctx context.Context, msgIds []string, userId int) ([]string, error) {
	var records []string

	err := t.Db.WithContext(ctx).Table("talk_records_delete").Select("msg_id").Where("user_id =?", userId).Where("msg_id in ?", msgIds).Scan(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
