package system

import (
	"context"
	"slices"

	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"github.com/samber/lo"
)

var _ admin.IMenuHandler = (*Menu)(nil)

type Menu struct {
	SysMenuRepo *repo.SysMenu
}

var tree = MenuTree{}

func (m *Menu) List(ctx context.Context, req *admin.MenuListRequest) (*admin.MenuListResponse, error) {
	items, err := m.SysMenuRepo.FindAll(ctx, func(db *gorm.DB) {
		db.Order("id asc")
	})

	if err != nil {
		return nil, err
	}

	return &admin.MenuListResponse{
		Items: tree.Build(items),
	}, nil
}

func (m *Menu) Create(ctx context.Context, in *admin.MenuCreateRequest) (*admin.MenuCreateResponse, error) {
	if in.ParentId > 0 {
		info, err := m.SysMenuRepo.FindById(ctx, in.ParentId)
		if err != nil {
			return nil, err
		}

		if in.MenuType == 3 && slices.Contains([]int32{1, 3}, info.MenuType) {
			return nil, errorx.New(400, "只能在页面菜单下添加按钮类型的子菜单")
		}
	} else {
		if in.MenuType == 3 {
			return nil, errorx.New(400, "只能在页面菜单下添加按钮类型的子菜单")
		}
	}

	data := &model.SysMenu{
		ParentId:  in.ParentId,
		Name:      in.Name,
		MenuType:  in.MenuType,
		Icon:      in.Icon,
		Path:      in.Path,
		Sort:      in.Sort,
		Hidden:    lo.Ternary(in.Hidden == "", "N", in.Hidden),
		UseLayout: lo.Ternary(in.UseLayout == "", "Y", in.UseLayout),
		AuthCode:  in.AuthCode,
		Status:    1,
	}

	err := m.SysMenuRepo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return &admin.MenuCreateResponse{Id: data.Id}, nil
}

func (m *Menu) Update(ctx context.Context, in *admin.MenuUpdateRequest) (*admin.MenuUpdateResponse, error) {
	_, err := m.SysMenuRepo.UpdateByWhere(ctx, map[string]any{
		"parent_id":  in.ParentId,
		"name":       in.Name,
		"icon":       in.Icon,
		"path":       in.Path,
		"sort":       in.Sort,
		"hidden":     lo.Ternary(in.Hidden == "", "N", in.Hidden),
		"status":     in.Status,
		"use_layout": lo.Ternary(in.UseLayout == "", "Y", in.UseLayout),
		"auth_code":  in.AuthCode,
	}, "id = ?", in.Id)
	if err != nil {
		return nil, err
	}

	return &admin.MenuUpdateResponse{Id: in.Id}, nil
}

func (m *Menu) Delete(ctx context.Context, in *admin.MenuDeleteRequest) (*admin.MenuDeleteResponse, error) {
	info, err := m.SysMenuRepo.FindById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if info.Status != 2 {
		return nil, errorx.New(400, "该菜单已启用，请先禁用后进行删除")
	}

	err = m.SysMenuRepo.Delete(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &admin.MenuDeleteResponse{Id: in.Id}, nil
}
