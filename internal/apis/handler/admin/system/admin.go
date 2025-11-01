package system

import (
	"context"
	"time"

	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

var _ admin.IAdminHandler = (*Admin)(nil)

type Admin struct {
	AdminRepo *repo.Admin
}

func (a *Admin) Create(ctx context.Context, in *admin.AdminCreateRequest) (*admin.AdminCreateResponse, error) {
	data := &model.Admin{
		Username:    in.Username,
		Password:    encrypt.HashPassword(in.Password),
		Gender:      3,
		Email:       in.Email,
		Status:      1,
		LastLoginAt: time.Now(),
	}

	err := a.AdminRepo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return &admin.AdminCreateResponse{Id: int32(data.Id)}, nil
}

func (a *Admin) List(ctx context.Context, in *admin.AdminListRequest) (*admin.AdminListResponse, error) {
	total, conditions, err := a.AdminRepo.Pagination(ctx, int(in.Page), int(in.PageSize), func(tx *gorm.DB) *gorm.DB {
		if in.Username != "" {
			tx = tx.Where("username = ?", in.Username)
		}

		if in.Email != "" {
			tx = tx.Where("email = ?", in.Email)
		}

		if in.Status > 0 {
			tx = tx.Where("status = ?", in.Status)
		}

		return tx.Order("id desc")
	})

	if err != nil {
		return nil, err
	}

	items := lo.Map(conditions, func(item *model.Admin, index int) *admin.AdminListResponse_Item {
		return &admin.AdminListResponse_Item{
			Id:          int32(item.Id),
			Username:    item.Username,
			Avatar:      item.Avatar,
			Mobile:      item.Mobile,
			Email:       item.Email,
			Status:      int32(item.Status),
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
			UpdatedAt:   item.UpdatedAt.Format(time.DateTime),
			LastLoginAt: item.LastLoginAt.Format(time.DateTime),
			RoleName:    "test",
		}
	})

	return &admin.AdminListResponse{
		Items:     items,
		Total:     int32(total),
		Page:      in.Page,
		PageSize:  in.PageSize,
		PageTotal: int32(total) / in.PageSize,
	}, nil
}

func (a *Admin) UpdateStatus(ctx context.Context, in *admin.AdminUpdateStatusRequest) (*admin.AdminUpdateStatusResponse, error) {
	_, err := a.AdminRepo.UpdateById(ctx, in.GetId(), map[string]any{
		"status": in.Status,
	})

	if err != nil {
		return nil, err
	}

	return &admin.AdminUpdateStatusResponse{Id: in.Id}, nil
}

func (a *Admin) ResetPassword(ctx context.Context, in *admin.AdminResetPasswordRequest) (*admin.AdminResetPasswordResponse, error) {
	_, err := a.AdminRepo.UpdateById(ctx, in.GetId(), map[string]any{
		"password": encrypt.HashPassword(in.Password),
	})

	if err != nil {
		return nil, err
	}

	return &admin.AdminResetPasswordResponse{Id: in.Id}, nil
}
