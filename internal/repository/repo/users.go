package repo

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Users struct {
	core.Repo[model.Users]
	tableCache core.TableCache[model.Users, int]
}

func NewUsers(db *gorm.DB, rds *redis.Client) *Users {
	return &Users{
		Repo:       core.NewRepo[model.Users](db),
		tableCache: core.NewTableCache[model.Users, int](rds),
	}
}

// Create 创建数据
func (u *Users) Create(user *model.Users) (*model.Users, error) {
	if err := u.Repo.Db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByMobile 手机号查询
func (u *Users) FindByMobile(ctx context.Context, mobile string) (*model.Users, error) {
	return u.Repo.FindByWhere(ctx, "mobile = ?", mobile)
}

// IsMobileExist 判断手机号是否存在
func (u *Users) IsMobileExist(ctx context.Context, mobile string) bool {
	exist, _ := u.Repo.IsExist(ctx, "mobile = ?", mobile)
	return exist
}

func (u *Users) FindByIdWithCache(ctx context.Context, id int) (*model.Users, error) {
	return u.tableCache.GetOrSet(ctx, id, func(ctx context.Context) (*model.Users, error) {
		return u.Repo.FindById(ctx, id)
	})
}

func (u *Users) ClearTableCache(ctx context.Context, id int) error {
	return u.tableCache.Del(ctx, id)
}
