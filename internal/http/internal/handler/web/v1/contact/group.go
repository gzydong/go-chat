package contact

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
)

type Group struct {
	service *service.ContactGroupService
	contact *service.ContactService
}

func NewGroup(service *service.ContactGroupService, contact *service.ContactService) *Group {
	return &Group{service: service, contact: contact}
}

// List 联系人分组列表
func (c *Group) List(ctx *ichat.Context) error {

	uid := ctx.UserId()

	items := make([]*web.ContactGroupListResponse_Item, 0)
	items = append(items, &web.ContactGroupListResponse_Item{
		Name:  "全部好友",
		Count: int32(len(c.contact.GetContactIds(ctx.Ctx(), uid))),
	})

	group, err := c.service.GetUserGroup(ctx.Ctx(), uid)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	for _, v := range group {
		items = append(items, &web.ContactGroupListResponse_Item{
			Id:    int32(v.Id),
			Name:  v.Name,
			Count: int32(v.Num),
		})
	}

	return ctx.Success(&web.ContactGroupListResponse{
		Items: items,
	})
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

	err := c.service.Repo().Create(ctx.Ctx(), data)
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

	affected, err := c.service.Repo().UpdateWhere(ctx.Ctx(), map[string]interface{}{
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

	err := c.service.Delete(ctx.Ctx(), int(params.Id), ctx.UserId())
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

	err := c.service.Sort(ctx.Ctx(), ctx.UserId(), items)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactGroupSortResponse{})
}
