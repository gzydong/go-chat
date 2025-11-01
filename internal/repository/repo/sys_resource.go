package repo

import (
	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SysResource struct {
	core.Repo[model.SysResource]
}

func NewSysResource(db *gorm.DB) *SysResource {
	return &SysResource{Repo: core.NewRepo[model.SysResource](db)}
}
