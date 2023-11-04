package repo

import (
	"context"
	"fmt"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Users struct {
	ichat.Repo[model.Users]
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{Repo: ichat.NewRepo[model.Users](db)}
}

// Create 创建数据
func (u *Users) Create(user *model.Users) (*model.Users, error) {
	if err := u.Repo.Db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByMobile 手机号查询
func (u *Users) FindByMobile(mobile string) (*model.Users, error) {

	if len(mobile) == 0 {
		return nil, fmt.Errorf("mobile is empty")
	}

	return u.Repo.FindByWhere(context.TODO(), "mobile = ?", mobile)
}

// IsMobileExist 判断手机号是否存在
func (u *Users) IsMobileExist(ctx context.Context, mobile string) bool {

	if len(mobile) == 0 {
		return false
	}

	exist, _ := u.Repo.QueryExist(ctx, "mobile = ?", mobile)
	return exist
}
