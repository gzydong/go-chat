package repo

import (
	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SysRole struct {
	core.Repo[model.SysRole]
}

func NewSysRole(db *gorm.DB) *SysRole {
	return &SysRole{Repo: core.NewRepo[model.SysRole](db)}
}
