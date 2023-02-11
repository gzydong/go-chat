package repo

import (
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ContactGroup struct {
	ichat.Repo[model.ContactGroup]
}

func NewContactGroup(db *gorm.DB) *ContactGroup {
	return &ContactGroup{Repo: ichat.NewRepo[model.ContactGroup](db)}
}
