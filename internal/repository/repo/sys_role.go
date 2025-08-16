package repo

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SysRole struct {
	core.Repo[model.SysRole]
}

func NewSysRole(db *gorm.DB) *SysRole {
	return &SysRole{Repo: core.NewRepo[model.SysRole](db)}
}
