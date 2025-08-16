package system

import (
	"time"

	"github.com/samber/lo"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type Role struct {
	SysRoleRepo *repo.SysRole
}

func (r *Role) List(ctx *core.Context) error {
	var in admin.RoleListRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	total, conditions, err := r.SysRoleRepo.Pagination(ctx.GetContext(), int(in.Page), int(in.PageSize), func(tx *gorm.DB) *gorm.DB {
		if in.RoleName != "" {
			tx = tx.Where("role_name = ?", in.RoleName)
		}

		if in.Status > 0 {
			tx = tx.Where("status = ?", in.Status)
		}

		return tx.Order("id desc")
	})

	if err != nil {
		return ctx.Error(err)
	}

	items := lo.Map(conditions, func(item *model.SysRole, index int) *admin.RoleListResponse_Item {
		return &admin.RoleListResponse_Item{
			Id:        int32(item.Id),
			RoleName:  item.RoleName,
			Status:    int32(item.Status),
			CreatedAt: item.CreatedAt.Format(time.DateTime),
			UpdatedAt: item.UpdatedAt.Format(time.DateTime),
			Explain:   item.Explain,
		}
	})

	return ctx.Success(&admin.RoleListResponse{
		Items:     items,
		Total:     int32(total),
		Page:      in.Page,
		PageSize:  in.PageSize,
		PageTotal: int32(total) / in.PageSize,
	})
}

func (r *Role) Create(ctx *core.Context) error {
	var in admin.RoleCreateRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	data := &model.SysRole{
		RoleName: in.RoleName,
		Explain:  in.Explain,
		Status:   1,
	}

	err := r.SysRoleRepo.Create(ctx.GetContext(), data)
	if err != nil {
		return err
	}

	return ctx.Success(admin.RoleCreateResponse{Id: int32(data.Id)})
}

func (r *Role) Update(ctx *core.Context) error {
	var in admin.RoleUpdateRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := r.SysRoleRepo.UpdateById(ctx.GetContext(), in.Id, map[string]any{
		"explain":   in.Explain,
		"role_name": in.RoleName,
	})

	if err != nil {
		return err
	}

	return ctx.Success(admin.RoleUpdateResponse{Id: in.Id})
}
