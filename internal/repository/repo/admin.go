package repo

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Admin struct {
	core.Repo[model.Admin]
}

func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{Repo: core.NewRepo[model.Admin](db)}
}
