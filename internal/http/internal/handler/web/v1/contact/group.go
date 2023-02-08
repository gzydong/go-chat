package contact

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
)

type Group struct {
	service *service.ContactGroupService
}

func NewGroup(service *service.ContactGroupService) *Group {
	return &Group{service: service}
}

func (c *Group) List(ctx *ichat.Context) error {
	return nil
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
