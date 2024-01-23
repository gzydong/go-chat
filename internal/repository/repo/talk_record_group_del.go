package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecordGroupDel struct {
	ichat.Repo[model.TalkRecordGroupDel]
}

func NewTalkRecordGroupDel(db *gorm.DB) *TalkRecordGroupDel {
	return &TalkRecordGroupDel{Repo: ichat.NewRepo[model.TalkRecordGroupDel](db)}
}

func (t *TalkRecordGroupDel) FindAllMsgIds(ctx context.Context, userId int, msgIds []string) ([]string, error) {
	var records []string
	err := t.Db.WithContext(ctx).Table("talk_record_group_del").Where("user_id = ? and msg_id in ?", userId, msgIds).Pluck("msg_id", &records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
