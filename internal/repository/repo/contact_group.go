package repo

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ContactGroup struct {
	core.Repo[model.ContactGroup]
}

func NewContactGroup(db *gorm.DB) *ContactGroup {
	return &ContactGroup{Repo: core.NewRepo[model.ContactGroup](db)}
}
