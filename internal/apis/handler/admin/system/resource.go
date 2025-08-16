package system

import (
	"time"

	"github.com/samber/lo"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type Resource struct {
	SysResourceRepo *repo.SysResource
}

func (a *Resource) List(ctx *core.Context) error {
	var in admin.ResourceListRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	total, conditions, err := a.SysResourceRepo.Pagination(ctx.GetContext(), int(in.Page), int(in.PageSize), func(tx *gorm.DB) *gorm.DB {
		if in.Uri != "" {
			// 这里使用模糊查询
			tx = tx.Where("uri like ?", in.Uri+"%")
		}

		if in.Name != "" {
			tx = tx.Where("name = ?", in.Name)
		}

		if in.Status > 0 {
			tx = tx.Where("status = ?", in.Status)
		}

		if in.Type > 0 {
			tx = tx.Where("`type` = ?", in.Type)
		}

		return tx.Order("id desc")
	})

	if err != nil {
		return ctx.Error(err)
	}

	items := lo.Map(conditions, func(item *model.SysResource, index int) *admin.ResourceListResponse_Item {
		return &admin.ResourceListResponse_Item{
			Id:        item.Id,
			Name:      item.Name,
			Uri:       item.Uri,
			Status:    item.Status,
			Type:      item.Type,
			CreatedAt: item.CreatedAt.Format(time.DateTime),
			UpdatedAt: item.UpdatedAt.Format(time.DateTime),
		}
	})

	return ctx.Success(&admin.ResourceListResponse{
		Items:     items,
		Total:     int32(total),
		Page:      in.Page,
		PageSize:  in.PageSize,
		PageTotal: int32(total) / in.PageSize,
	})
}

func (a *Resource) Create(ctx *core.Context) error {
	var in admin.ResourceCreateRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	data := &model.SysResource{
		Name:   in.Name,
		Uri:    in.Uri,
		Type:   in.Type,
		Status: 1,
	}

	err := a.SysResourceRepo.Create(ctx.GetContext(), data)
	if err != nil {
		return err
	}

	return ctx.Success(admin.AdminCreateResponse{Id: data.Id})
}

func (a *Resource) Update(ctx *core.Context) error {
	var in admin.ResourceUpdateRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := a.SysResourceRepo.UpdateById(ctx.GetContext(), in.GetId(), map[string]any{
		"status": in.Status,
		"name":   in.Name,
		"uri":    in.Uri,
	})

	if err != nil {
		return err
	}

	return ctx.Success(admin.ResourceUpdateResponse{Id: in.Id})
}

func (a *Resource) Delete(ctx *core.Context) error {
	var in admin.ResourceDeleteRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	resource, err := a.SysResourceRepo.FindById(ctx.GetContext(), in.Id)
	if err != nil {
		return err
	}

	if resource.Status == model.ResourceStatusNormal {
		return errorx.New(400, "该资源已启用，请先禁用后进行删除")
	}

	err = a.SysResourceRepo.Delete(ctx.GetContext(), resource.Id)
	if err != nil {
		return err
	}

	return ctx.Success(admin.ResourceDeleteResponse{Id: in.Id})
}
