package system

import (
	"context"
	"time"

	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

var _ admin.IResourceHandler = (*Resource)(nil)

type Resource struct {
	SysResourceRepo *repo.SysResource
}

func (a *Resource) List(ctx context.Context, in *admin.ResourceListRequest) (*admin.ResourceListResponse, error) {
	total, conditions, err := a.SysResourceRepo.Pagination(ctx, int(in.Page), int(in.PageSize), func(tx *gorm.DB) *gorm.DB {
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
		return nil, err
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

	return &admin.ResourceListResponse{
		Items:     items,
		Total:     int32(total),
		Page:      in.Page,
		PageSize:  in.PageSize,
		PageTotal: int32(total) / in.PageSize,
	}, nil
}

func (a *Resource) Create(ctx context.Context, in *admin.ResourceCreateRequest) (*admin.ResourceCreateResponse, error) {
	data := &model.SysResource{
		Name:   in.Name,
		Uri:    in.Uri,
		Type:   in.Type,
		Status: 1,
	}

	err := a.SysResourceRepo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return &admin.ResourceCreateResponse{Id: data.Id}, nil
}

func (a *Resource) Update(ctx context.Context, in *admin.ResourceUpdateRequest) (*admin.ResourceUpdateResponse, error) {
	_, err := a.SysResourceRepo.UpdateById(ctx, in.GetId(), map[string]any{
		"status": in.Status,
		"name":   in.Name,
		"uri":    in.Uri,
	})

	if err != nil {
		return nil, err
	}

	return &admin.ResourceUpdateResponse{Id: in.Id}, nil
}

func (a *Resource) Delete(ctx context.Context, in *admin.ResourceDeleteRequest) (*admin.ResourceDeleteResponse, error) {
	resource, err := a.SysResourceRepo.FindById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if resource.Status == model.ResourceStatusNormal {
		return nil, errorx.New(400, "该资源已启用，请先禁用后进行删除")
	}

	err = a.SysResourceRepo.Delete(ctx, resource.Id)
	if err != nil {
		return nil, err
	}

	return &admin.ResourceDeleteResponse{Id: in.Id}, nil
}
