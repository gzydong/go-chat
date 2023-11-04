package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Robot struct {
	ichat.Repo[model.Robot]
}

func NewRobot(db *gorm.DB) *Robot {
	return &Robot{Repo: ichat.NewRepo[model.Robot](db)}
}

// GetLoginRobot 获取登录机器的信息
func (r *Robot) GetLoginRobot(ctx context.Context) (*model.Robot, error) {
	return r.Repo.FindByWhere(ctx, "type = ? and status = ?", 1, model.RootStatusNormal)
}
