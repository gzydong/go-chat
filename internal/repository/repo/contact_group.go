package repo

import (
	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ContactGroup struct {
	core.Repo[model.ContactGroup]
}

func NewContactGroup(db *gorm.DB) *ContactGroup {
	return &ContactGroup{Repo: core.NewRepo[model.ContactGroup](db)}
}
