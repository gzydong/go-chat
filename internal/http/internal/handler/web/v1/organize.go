package v1

import (
	"strconv"
	"strings"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"go-chat/internal/service/organize"
)

type Organize struct {
	deptServ     *organize.DeptService
	organizeServ *organize.OrganizeService
	positionServ *organize.PositionService
}

func NewOrganize(deptServ *organize.DeptService, organizeServ *organize.OrganizeService, positionServ *organize.PositionService) *Organize {
	return &Organize{deptServ: deptServ, organizeServ: organizeServ, positionServ: positionServ}
}

func (o *Organize) DepartmentList(ctx *ichat.Context) error {

	uid := ctx.UserId()
	if isOk, _ := o.organizeServ.Dao().IsQiyeMember(uid); !isOk {
		return ctx.Success(entity.H{
			"items": []string{},
		})
	}

	items, err := o.deptServ.Dao().FindAll()
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(entity.H{
		"items": items,
	})
}

type UserInfo struct {
	UserId        int              `json:"user_id"`
	Nickname      string           `json:"nickname"`
	Gender        int              `json:"gender"`
	PositionItems []*PositionItems `json:"position_items" gorm:"-"`
	DeptItems     []*DeptItems     `json:"dept_items" gorm:"-"`
}

type DeptItems struct {
	DeptId    int    `json:"dept_id"`
	DeptName  string `json:"dept_name"`
	Ancestors string `json:"ancestors"`
}

type PositionItems struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sort int    `json:"sort"`
}

func (o *Organize) PersonnelList(ctx *ichat.Context) error {

	uid := ctx.UserId()
	if isOk, _ := o.organizeServ.Dao().IsQiyeMember(uid); !isOk {
		return ctx.Success(entity.H{
			"items": []string{},
		})
	}

	list, err := o.organizeServ.Dao().FindAll()
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	departments, err := o.deptServ.Dao().FindAll()
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	deptHash := make(map[int]*model.OrganizeDept)
	for _, department := range departments {
		deptHash[department.DeptId] = department
	}

	positions, err := o.positionServ.Dao().FindAll()
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	positionHash := make(map[int]*model.OrganizePost)
	for _, position := range positions {
		positionHash[position.PositionId] = position
	}

	items := make([]*UserInfo, 0)
	for _, info := range list {
		data := &UserInfo{
			UserId:        info.UserId,
			Nickname:      info.Nickname,
			Gender:        info.Gender,
			PositionItems: make([]*PositionItems, 0),
			DeptItems:     make([]*DeptItems, 0),
		}

		for _, key := range strings.Split(info.Department, ",") {
			id, _ := strconv.Atoi(key)
			if val, ok := deptHash[id]; ok {
				data.DeptItems = append(data.DeptItems, &DeptItems{
					DeptId:    val.DeptId,
					DeptName:  val.DeptName,
					Ancestors: val.Ancestors,
				})
			}
		}

		for _, key := range strings.Split(info.Position, ",") {
			id, _ := strconv.Atoi(key)
			if val, ok := positionHash[id]; ok {
				data.PositionItems = append(data.PositionItems, &PositionItems{
					Code: val.PostCode,
					Name: val.PostName,
					Sort: val.Sort,
				})
			}
		}

		items = append(items, data)
	}

	return ctx.Success(entity.H{
		"items": items,
	})
}
