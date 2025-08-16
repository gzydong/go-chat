package repo

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SysResource struct {
	core.Repo[model.SysResource]
}

func NewSysResource(db *gorm.DB) *SysResource {
	return &SysResource{Repo: core.NewRepo[model.SysResource](db)}
}
