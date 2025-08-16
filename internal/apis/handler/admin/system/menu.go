package system

import (
	"slices"
	"sort"

	"github.com/samber/lo"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type Menu struct {
	SysMenuRepo *repo.SysMenu
}

func (m *Menu) List(ctx *core.Context) error {
	items, err := m.SysMenuRepo.FindAll(ctx.GetContext(), func(db *gorm.DB) {
		db.Order("id asc")
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(map[string]any{
		"items": m.SysMenuRepo.BuildMenuTree(items),
	})
}

func (m *Menu) Create(ctx *core.Context) error {
	var in admin.MenuCreateRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	if in.ParentId > 0 {
		info, err := m.SysMenuRepo.FindById(ctx.GetContext(), in.ParentId)
		if err != nil {
			return err
		}

		if in.MenuType == 3 && slices.Contains([]int32{1, 3}, info.MenuType) {
			return ctx.InvalidParams("只能在页面菜单下添加按钮类型的子菜单")
		}
	} else {
		if in.MenuType == 3 {
			return ctx.InvalidParams("只能在页面菜单下添加按钮类型的子菜单")
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

	err := m.SysMenuRepo.Create(ctx.GetContext(), data)
	if err != nil {
		return err
	}

	return ctx.Success(admin.MenuCreateResponse{Id: data.Id})
}

func (m *Menu) Update(ctx *core.Context) error {
	var in admin.MenuUpdateRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := m.SysMenuRepo.UpdateByWhere(ctx.GetContext(), map[string]any{
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
		return err
	}

	return ctx.Success(admin.MenuCreateResponse{Id: in.Id})
}

func (m *Menu) Delete(ctx *core.Context) error {
	var in admin.MenuDeleteRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	info, err := m.SysMenuRepo.FindById(ctx.GetContext(), in.Id)
	if err != nil {
		return err
	}

	if info.Status != 2 {
		return errorx.New(400, "该菜单已启用，请先禁用后进行删除")
	}

	err = m.SysMenuRepo.Delete(ctx.GetContext(), in.Id)
	if err != nil {
		return err
	}

	return ctx.Success(admin.MenuDeleteResponse{Id: in.Id})
}

type UserMenus struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Meta struct {
		Icon      string `json:"icon"`
		Title     string `json:"title"`
		Sort      int32  `json:"sort"`
		Hidden    string `json:"hidden"`
		UseLayout string `json:"use_layout"`
		FrameSrc  string `json:"frame_src,omitempty"`
	} `json:"meta"`
	Children []UserMenus `json:"children"`
}

// buildUserMenus 递归构建UserMenus结构
func (m *Menu) buildUserMenus(menuItems []*repo.MenuItem) []UserMenus {
	var userMenus []UserMenus

	for _, item := range menuItems {
		if item.Status != 1 { // 假设1为启用状态
			continue
		}

		userMenu := UserMenus{
			Path: item.Path,
			Name: item.Name,
		}

		// 设置Meta信息
		userMenu.Meta.Icon = item.Icon
		userMenu.Meta.Title = item.Name
		userMenu.Meta.Sort = item.Sort
		userMenu.Meta.Hidden = item.Hidden
		userMenu.Meta.UseLayout = item.UseLayout
		userMenu.Meta.FrameSrc = ""

		// 如果有子菜单，递归处理
		if len(item.Children) > 0 {
			userMenu.Children = m.buildUserMenus(item.Children)
		}

		userMenus = append(userMenus, userMenu)
	}

	sort.Slice(userMenus, func(i, j int) bool {
		return userMenus[i].Meta.Sort < userMenus[j].Meta.Sort
	})

	return userMenus
}

// GetUserMenus 获取用户菜单列表
func (m *Menu) GetUserMenus(ctx *core.Context) error {
	items, err := m.SysMenuRepo.FindAll(ctx.GetContext(), func(db *gorm.DB) {
		db.Where("status = ?", 1)
		db.Where("menu_type in ?", []int32{1, 2})
		db.Order("id asc")
	})
	if err != nil {
		return ctx.Error(err)
	}

	userMenus := m.buildUserMenus(m.SysMenuRepo.BuildMenuTree(items))
	return ctx.Success(map[string]any{
		"items": userMenus,
	})
}
