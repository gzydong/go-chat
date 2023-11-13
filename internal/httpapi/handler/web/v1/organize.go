package v1

import (
	"strconv"
	"strings"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type Organize struct {
	DepartmentRepo *repo.Department
	PositionRepo   *repo.Position
	OrganizeRepo   *repo.Organize
}

func (o *Organize) DepartmentList(ctx *ichat.Context) error {

	uid := ctx.UserId()
	if isOk, _ := o.OrganizeRepo.IsQiyeMember(ctx.Ctx(), uid); !isOk {
		return ctx.Success(&web.OrganizeDepartmentListResponse{})
	}

	list, err := o.DepartmentRepo.List(ctx.Ctx())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	items := make([]*web.OrganizeDepartmentListResponse_Item, 0, len(list))
	for _, dept := range list {
		items = append(items, &web.OrganizeDepartmentListResponse_Item{
			DeptId:    int32(dept.DeptId),
			ParentId:  int32(dept.ParentId),
			DeptName:  dept.DeptName,
			Ancestors: dept.Ancestors,
		})
	}

	return ctx.Success(&web.OrganizeDepartmentListResponse{Items: items})
}

func (o *Organize) PersonnelList(ctx *ichat.Context) error {

	uid := ctx.UserId()
	if isOk, _ := o.OrganizeRepo.IsQiyeMember(ctx.Ctx(), uid); !isOk {
		return ctx.Success(&web.OrganizePersonnelListResponse{})
	}

	list, err := o.OrganizeRepo.List()
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	departments, err := o.DepartmentRepo.List(ctx.Ctx())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	deptHash := make(map[int]*model.OrganizeDept)
	for _, department := range departments {
		deptHash[department.DeptId] = department
	}

	positions, err := o.PositionRepo.List(ctx.Ctx())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
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
			Gender:        int32(info.Gender),
			PositionItems: make([]*web.OrganizePersonnelListResponse_Position, 0),
			DeptItems:     make([]*web.OrganizePersonnelListResponse_Dept, 0),
		}

		for _, key := range strings.Split(info.Department, ",") {
			id, _ := strconv.Atoi(key)
			if val, ok := deptHash[id]; ok {
				data.DeptItems = append(data.DeptItems, &web.OrganizePersonnelListResponse_Dept{
					DeptId:    int32(val.DeptId),
					DeptName:  val.DeptName,
					Ancestors: val.Ancestors,
				})
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
