package system

import (
	"sort"
	"time"

	"github.com/gzydong/go-chat/api/pb/admin/v1"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/samber/lo"
)

type MenuTree struct {
}

// MenuItem 菜单节点结构
type MenuItem struct {
	Id        int32       `json:"id"`
	ParentId  int32       `json:"parent_id"`
	Name      string      `json:"name"`
	MenuType  int32       `json:"menu_type"`
	Icon      string      `json:"icon"`
	Path      string      `json:"path"`
	Sort      int32       `json:"sort"`
	Status    int32       `json:"status"`
	Hidden    string      `json:"hidden"`
	UseLayout string      `json:"use_layout"`
	AuthCode  string      `json:"auth_code"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Children  []*MenuItem `json:"children"`
}

// sortMenuItems 递归排序菜单项
func sortMenuItems(items []*admin.MenuItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Sort < items[j].Sort
	})

	for _, item := range items {
		if len(item.Children) > 0 {
			sortMenuItems(item.Children)
		}
	}
}

func (m *MenuTree) Build(items []*model.SysMenu) []*admin.MenuItem {
	list := lo.Map(items, func(item *model.SysMenu, index int) *admin.MenuItem {
		return &admin.MenuItem{
			Id:        item.Id,
			ParentId:  item.ParentId,
			Name:      item.Name,
			MenuType:  item.MenuType,
			Icon:      item.Icon,
			Path:      item.Path,
			Sort:      item.Sort,
			Status:    item.Status,
			Hidden:    item.Hidden,
			UseLayout: item.UseLayout,
			AuthCode:  item.AuthCode,
			CreatedAt: item.CreatedAt.Format(time.DateTime),
			UpdatedAt: item.UpdatedAt.Format(time.DateTime),
			Children:  nil,
		}
	})

	return m.toTree(list)
}

func (m *MenuTree) toTree(menuItems []*admin.MenuItem) []*admin.MenuItem {
	menuMap := make(map[int32]*admin.MenuItem)
	var rootMenus []*admin.MenuItem

	// 初始化所有节点
	for _, item := range menuItems {
		item.Children = make([]*admin.MenuItem, 0)
		menuMap[item.Id] = item
	}

	// 构建父子关系
	for _, item := range menuItems {
		if item.ParentId == 0 {
			rootMenus = append(rootMenus, item)
		} else {
			if parent, exists := menuMap[item.ParentId]; exists {
				parent.Children = append(parent.Children, item)
			}
		}
	}

	// 排序
	sortMenuItems(rootMenus)

	return rootMenus
}
