package repo

import (
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Admin struct {
	ichat.Repo[model.Admin]
}

func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{Repo: ichat.NewRepo[model.Admin](db)}
}
