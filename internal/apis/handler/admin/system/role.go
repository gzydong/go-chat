package system

import (
	"context"
	"time"

	"go-chat/api/pb/admin/v1"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

var _ admin.IRoleHandler = (*Role)(nil)

type Role struct {
	SysRoleRepo *repo.SysRole
}

func (r *Role) List(ctx context.Context, in *admin.RoleListRequest) (*admin.RoleListResponse, error) {
	total, conditions, err := r.SysRoleRepo.Pagination(ctx, int(in.Page), int(in.PageSize), func(tx *gorm.DB) *gorm.DB {
		if in.RoleName != "" {
			tx = tx.Where("role_name = ?", in.RoleName)
		}

		if in.Status > 0 {
			tx = tx.Where("status = ?", in.Status)
		}

		return tx.Order("id desc")
	})

	if err != nil {
		return nil, err
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

	return &admin.RoleListResponse{
		Items:     items,
		Total:     int32(total),
		Page:      in.Page,
		PageSize:  in.PageSize,
		PageTotal: int32(total) / in.PageSize,
	}, nil
}

func (r *Role) Create(ctx context.Context, in *admin.RoleCreateRequest) (*admin.RoleCreateResponse, error) {
	data := &model.SysRole{
		RoleName: in.RoleName,
		Explain:  in.Explain,
		Status:   1,
	}

	err := r.SysRoleRepo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return &admin.RoleCreateResponse{Id: int32(data.Id)}, nil
}

func (r *Role) Update(ctx context.Context, in *admin.RoleUpdateRequest) (*admin.RoleUpdateResponse, error) {
	_, err := r.SysRoleRepo.UpdateById(ctx, in.Id, map[string]any{
		"explain":   in.Explain,
		"role_name": in.RoleName,
	})

	if err != nil {
		return nil, err
	}

	return &admin.RoleUpdateResponse{Id: in.Id}, nil
}
