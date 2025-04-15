package v1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type Organize struct {
	DepartmentRepo *repo.Department
	PositionRepo   *repo.Position
	OrganizeRepo   *repo.Organize
}

func (o *Organize) DepartmentList(ctx *core.Context) error {
	uid := ctx.AuthId()
	if isOk, _ := o.OrganizeRepo.IsQiyeMember(ctx.GetContext(), uid); !isOk {
		return ctx.Success(&web.OrganizeDepartmentListResponse{})
	}

	list, err := o.DepartmentRepo.List(ctx.GetContext())
	if err != nil {
		return ctx.Error(err)
	}

	// 部门分组统计
	groups, err := o.OrganizeRepo.DepartmentGroupCount(ctx.GetContext())
	if err != nil {
		return ctx.Error(err)
	}

	groupsHash := make(map[int32]int32)
	for _, v := range groups {
		groupsHash[v.DeptId] = v.Count
	}

	var mapping = make(map[string]int32)
	for _, v := range list {
		mapping[fmt.Sprintf("%s,%d", v.Ancestors, v.DeptId)] = groupsHash[int32(v.DeptId)]
	}

	items := make([]*web.OrganizeDepartmentListResponse_Item, 0, len(list))
	items = append(items, &web.OrganizeDepartmentListResponse_Item{
		DeptId:    -1,
		ParentId:  0,
		DeptName:  "企业成员",
		Ancestors: "",
		Count:     lo.SumBy(groups, func(item *repo.GroupCount) int32 { return item.Count }),
	})

	for _, dept := range list {
		var count int32 = 0

		s := fmt.Sprintf("%s,%d", dept.Ancestors, dept.DeptId)
		for key, value := range mapping {
			if strings.HasPrefix(key, s) {
				count += value
			}
		}

		items = append(items, &web.OrganizeDepartmentListResponse_Item{
			DeptId:    int32(dept.DeptId),
			ParentId:  int32(dept.ParentId),
			DeptName:  dept.DeptName,
			Ancestors: dept.Ancestors,
			Count:     count,
		})
	}

	return ctx.Success(&web.OrganizeDepartmentListResponse{Items: items})
}

func (o *Organize) PersonnelList(ctx *core.Context) error {

	// 判断是否是企业成员
	uid := ctx.AuthId()
	if isOk, _ := o.OrganizeRepo.IsQiyeMember(ctx.GetContext(), uid); !isOk {
		return ctx.Success(&web.OrganizePersonnelListResponse{})
	}

	list, err := o.OrganizeRepo.List()
	if err != nil {
		return ctx.Error(err)
	}

	departments, err := o.DepartmentRepo.List(ctx.GetContext())
	if err != nil {
		return ctx.Error(err)
	}

	deptHash := make(map[int]*model.OrganizeDept)
	for _, department := range departments {
		deptHash[department.DeptId] = department
	}

	positions, err := o.PositionRepo.List(ctx.GetContext())
	if err != nil {
		return ctx.Error(err)
	}

	positionHash := make(map[int]*model.OrganizePost)
	for _, position := range positions {
		positionHash[position.PositionId] = position
	}

	items := make([]*web.OrganizePersonnelListResponse_Item, 0)
	for _, info := range list {
		data := &web.OrganizePersonnelListResponse_Item{
			UserId:        int32(info.UserId),
			Nickname:      info.Nickname,
			Avatar:        info.Avatar,
			Gender:        int32(info.Gender),
			PositionItems: make([]*web.OrganizePersonnelListResponse_Position, 0),
			DeptItem:      &web.OrganizePersonnelListResponse_Dept{},
		}

		// 目前仅支持一个人一个部门
		if val, ok := deptHash[info.Department]; ok {
			data.DeptItem = &web.OrganizePersonnelListResponse_Dept{
				DeptId:    int32(info.Department),
				DeptName:  val.DeptName,
				Ancestors: val.Ancestors,
			}
		}

		for _, key := range strings.Split(info.Position, ",") {
			id, _ := strconv.Atoi(key)
			if val, ok := positionHash[id]; ok {
				data.PositionItems = append(data.PositionItems, &web.OrganizePersonnelListResponse_Position{
					Code: val.PostCode,
					Name: val.PostName,
					Sort: int32(val.Sort),
				})
			}
		}

		items = append(items, data)
	}

	return ctx.Success(&web.OrganizePersonnelListResponse{Items: items})
}
