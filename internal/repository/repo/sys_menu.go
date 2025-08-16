package repo

import (
	"sort"
	"time"

	"github.com/samber/lo"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SysMenu struct {
	core.Repo[model.SysMenu]
}

func NewSysMenu(db *gorm.DB) *SysMenu {
	return &SysMenu{Repo: core.NewRepo[model.SysMenu](db)}
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

// BuildMenuTree 构建菜单树
func BuildMenuTree(menuItems []*MenuItem) []*MenuItem {
	menuMap := make(map[int32]*MenuItem)
	var rootMenus []*MenuItem

	// 初始化所有节点
	for _, item := range menuItems {
		item.Children = []*MenuItem{}
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

// sortMenuItems 递归排序菜单项
func sortMenuItems(items []*MenuItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Sort < items[j].Sort
	})

	for _, item := range items {
		if len(item.Children) > 0 {
			sortMenuItems(item.Children)
		}
	}
}

func (s *SysMenu) BuildMenuTree(items []*model.SysMenu) []*MenuItem {
	list := lo.Map(items, func(item *model.SysMenu, index int) *MenuItem {
		return &MenuItem{
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

	return BuildMenuTree(list)
}
