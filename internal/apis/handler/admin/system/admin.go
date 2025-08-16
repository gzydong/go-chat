package system

import (
	"time"

	"github.com/samber/lo"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type Admin struct {
	AdminRepo *repo.Admin
}

func (a *Admin) List(ctx *core.Context) error {
	var in admin.AdminListRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	total, conditions, err := a.AdminRepo.Pagination(ctx.GetContext(), int(in.Page), int(in.PageSize), func(tx *gorm.DB) *gorm.DB {
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
		return ctx.Error(err)
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

	return ctx.Success(&admin.AdminListResponse{
		Items:     items,
		Total:     int32(total),
		Page:      in.Page,
		PageSize:  in.PageSize,
		PageTotal: int32(total) / in.PageSize,
	})
}

func (a *Admin) Create(ctx *core.Context) error {
	var in admin.AdminCreateRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	data := &model.Admin{
		Username:    in.Username,
		Password:    encrypt.HashPassword(in.Password),
		Gender:      3,
		Email:       in.Email,
		Status:      1,
		LastLoginAt: time.Now(),
	}

	err := a.AdminRepo.Create(ctx.GetContext(), data)
	if err != nil {
		return err
	}

	return ctx.Success(admin.AdminCreateResponse{Id: int32(data.Id)})
}

func (a *Admin) UpdateStatus(ctx *core.Context) error {
	var in admin.AdminStatusRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := a.AdminRepo.UpdateById(ctx.GetContext(), in.GetId(), map[string]any{
		"status": in.Status,
	})

	if err != nil {
		return err
	}

	return ctx.Success(admin.AdminStatusResponse{Id: in.Id})
}

func (a *Admin) ResetPassword(ctx *core.Context) error {
	var in admin.AdminResetPasswordRequest

	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := a.AdminRepo.UpdateById(ctx.GetContext(), in.GetId(), map[string]any{
		"password": encrypt.HashPassword(in.Password),
	})

	if err != nil {
		return err
	}

	return ctx.Success(admin.AdminResetPasswordResponse{Id: in.Id})
}
