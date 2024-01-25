package repo

import (
	"context"

	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkGroupMessageDel struct {
	core.Repo[model.TalkGroupMessageDel]
}

func NewTalkRecordGroupDel(db *gorm.DB) *TalkGroupMessageDel {
	return &TalkGroupMessageDel{Repo: core.NewRepo[model.TalkGroupMessageDel](db)}
}

func (t *TalkGroupMessageDel) FindAllMsgIds(ctx context.Context, userId int, msgIds []string) ([]string, error) {
	var records []string
	err := t.Db.WithContext(ctx).Table("talk_group_message_del").Where("user_id = ? and msg_id in ?", userId, msgIds).Pluck("msg_id", &records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
