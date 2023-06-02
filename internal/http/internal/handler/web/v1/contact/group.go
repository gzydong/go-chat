package contact

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
	"gorm.io/gorm"
)

type Group struct {
	contactGroupService *service.ContactGroupService
	contactService      *service.ContactService
}

func NewGroup(contactGroupService *service.ContactGroupService, contactService *service.ContactService) *Group {
	return &Group{contactGroupService: contactGroupService, contactService: contactService}
}

// List 联系人分组列表
func (c *Group) List(ctx *ichat.Context) error {

	uid := ctx.UserId()

	items := make([]*web.ContactGroupListResponse_Item, 0)
	items = append(items, &web.ContactGroupListResponse_Item{
		Name:  "全部好友",
		Count: int32(len(c.contactService.GetContactIds(ctx.Ctx(), uid))),
	})

	group, err := c.contactGroupService.GetUserGroup(ctx.Ctx(), uid)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	for _, v := range group {
		items = append(items, &web.ContactGroupListResponse_Item{
			Id:    int32(v.Id),
			Name:  v.Name,
			Count: int32(v.Num),
			Sort:  int32(v.Sort),
		})
	}

	return ctx.Success(&web.ContactGroupListResponse{Items: items})
}

func (c *Group) Create(ctx *ichat.Context) error {

	params := &web.ContactGroupCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	data := &model.ContactGroup{
		UserId: ctx.UserId(),
		Name:   params.GetName(),
		Sort:   int(params.GetSort()),
	}

	err := c.contactGroupService.Repo().Create(ctx.Ctx(), data)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactGroupCreateResponse{Id: int32(data.Id)})
}

func (c *Group) Update(ctx *ichat.Context) error {

	params := &web.ContactGroupUpdateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	affected, err := c.contactGroupService.Repo().UpdateWhere(ctx.Ctx(), map[string]any{
		"name":       params.Name,
		"sort":       params.Sort,
		"updated_at": timeutil.DateTime(),
	}, "id = ? and user_id = ?", params.Id, ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if affected == 0 {
		return ctx.ErrorBusiness("数据不存在")
	}

	return ctx.Success(&web.ContactGroupUpdateResponse{Id: params.Id})
}

func (c *Group) Delete(ctx *ichat.Context) error {

	params := &web.ContactGroupDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.contactGroupService.Delete(ctx.Ctx(), int(params.Id), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactGroupDeleteResponse{Id: params.Id})
}

func (c *Group) Sort(ctx *ichat.Context) error {
	params := &web.ContactGroupSortRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	items := make([]*model.ContactGroup, 0, len(params.GetItems()))
	for _, item := range params.GetItems() {
		items = append(items, &model.ContactGroup{
			Id:   int(item.Id),
			Sort: int(item.Sort),
		})
	}

	err := c.contactGroupService.Sort(ctx.Ctx(), ctx.UserId(), items)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactGroupSortResponse{})
}

func (c *Group) Save(ctx *ichat.Context) error {
	params := &web.ContactGroupSaveRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	updateItems := make([]*model.ContactGroup, 0)
	deleteItems := make([]int, 0)
	insertItems := make([]*model.ContactGroup, 0)

	ids := make(map[int]struct{})
	for i, item := range params.GetItems() {
		if item.Id > 0 {
			ids[int(item.Id)] = struct{}{}
			updateItems = append(updateItems, &model.ContactGroup{
				Id:   int(item.Id),
				Sort: i + 1,
				Name: item.Name,
			})
		} else {
			insertItems = append(insertItems, &model.ContactGroup{
				Sort:   i + 1,
				Name:   item.Name,
				UserId: uid,
			})
		}
	}

	all, err := c.contactGroupService.Repo().FindAll(ctx.Ctx())
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	for _, m := range all {
		if _, ok := ids[m.Id]; !ok {
			deleteItems = append(deleteItems, m.Id)
		}
	}

	err = c.contactGroupService.Db().Transaction(func(tx *gorm.DB) error {

		if len(insertItems) > 0 {
			if err := tx.Create(insertItems).Error; err != nil {
				return err
			}
		}

		if len(deleteItems) > 0 {
			err := tx.Delete(model.ContactGroup{}, "id in (?) and user_id = ?", deleteItems, uid).Error
			if err != nil {
				return err
			}

			tx.Table("contact").
				Where("user_id = ? and group_id in (?)", uid, deleteItems).
				UpdateColumn("group_id", 0)
		}

		for _, item := range updateItems {
			err = tx.Table("contact_group").
				Where("id = ? and user_id = ?", item.Id, uid).
				Updates(map[string]any{
					"name": item.Name,
					"sort": item.Sort,
				}).Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(&web.ContactGroupSaveResponse{})
}
