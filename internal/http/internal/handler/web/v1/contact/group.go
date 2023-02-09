package contact

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
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

func (c *Group) List(ctx *ichat.Context) error {

	items := make([]*web.ContactGroupListResponse_Item, 0)

	items = append(items, &web.ContactGroupListResponse_Item{
		Id:    0,
		Name:  "全部好友",
		Count: int32(len(c.contact.GetContactIds(ctx.Ctx(), ctx.UserId()))),
	})

	items = append(items, &web.ContactGroupListResponse_Item{
		Id:    1,
		Name:  "同事",
		Count: 0,
	})

	items = append(items, &web.ContactGroupListResponse_Item{
		Id:    2,
		Name:  "朋友",
		Count: 0,
	})

	items = append(items, &web.ContactGroupListResponse_Item{
		Id:    3,
		Name:  "家人",
		Count: 0,
	})

	items = append(items, &web.ContactGroupListResponse_Item{
		Id:    4,
		Name:  "陌生人",
		Count: 0,
	})

	return ctx.Success(&web.ContactGroupListResponse{
		Items: items,
	})
}

func (c *Group) Create(ctx *ichat.Context) error {

	params := &web.ContactGroupCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	id, err := c.service.Create(ctx.Ctx(), &model.ContactGroup{
		UserId: ctx.UserId(),
		Name:   params.GetName(),
		Sort:   int(params.GetSort()),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactGroupCreateResponse{Id: int32(id)})
}

func (c *Group) Update(ctx *ichat.Context) error {
	return nil
}

func (c *Group) Delete(ctx *ichat.Context) error {
	return nil
}

func (c *Group) Sort(ctx *ichat.Context) error {
	return nil
}
