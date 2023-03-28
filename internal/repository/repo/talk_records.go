package repo

import (
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecords struct {
	ichat.Repo[model.TalkRecords]
}

func NewTalkRecords(db *gorm.DB) *TalkRecords {
	return &TalkRecords{Repo: ichat.NewRepo[model.TalkRecords](db)}
}
