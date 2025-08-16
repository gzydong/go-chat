package repo

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SysAdminTotp struct {
	core.Repo[model.SysAdminTotp]
}

func NewSysAdminTotp(db *gorm.DB) *SysAdminTotp {
	return &SysAdminTotp{Repo: core.NewRepo[model.SysAdminTotp](db)}
}
