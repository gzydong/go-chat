package user

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ admin.IUserHandler = (*User)(nil)

type User struct {
	UserRepo *repo.Users
}

func (u *User) List(ctx context.Context, req *admin.UserListRequest) (*admin.UserListResponse, error) {
	total, data, err := u.UserRepo.Pagination(ctx, int(req.Page), int(req.PageSize), func(tx *gorm.DB) *gorm.DB {
		if len(req.Mobile) > 0 {
			tx = tx.Where("mobile = ?", req.Mobile)
		}

		if len(req.Email) > 0 {
			tx = tx.Where("email = ?", req.Email)
		}

		if req.Status > 0 {
			tx = tx.Where("status = ?", req.Status)
		}

		tx = tx.Order("id desc")
		return tx
	})

	if err != nil {
		return nil, err
	}

	pageTotal := 0

	if total > 0 {
		pageTotal = int(total) / int(req.PageSize)
	}

	return &admin.UserListResponse{
		Items: lo.Map(data, func(item *model.Users, index int) *admin.UserListResponse_Item {
			return &admin.UserListResponse_Item{
				Id:          int32(item.Id),
				Username:    item.Nickname,
				Email:       item.Email,
				Mobile:      lo.FromPtr(item.Mobile),
				Status:      int32(item.Status),
				CreatedAt:   item.CreatedAt.Format(time.DateTime),
				UpdatedAt:   item.UpdatedAt.Format(time.DateTime),
				LastLoginAt: item.UpdatedAt.Format(time.DateTime),
				Avatar:      item.Avatar,
			}
		}),
		Total:     int32(total),
		Page:      req.Page,
		PageSize:  req.PageSize,
		PageTotal: int32(pageTotal),
	}, nil
}

func (u *User) Update(ctx context.Context, req *admin.UserUpdateRequest) (*admin.UserUpdateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *User) Detail(ctx context.Context, req *admin.UserDetailRequest) (*admin.UserDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}
