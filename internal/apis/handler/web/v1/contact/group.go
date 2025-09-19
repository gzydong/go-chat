package contact

import (
	"context"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"gorm.io/gorm"
)

var _ web.IContactGroupHandler = (*Group)(nil)

type Group struct {
	ContactRepo         *repo.Contact
	ContactGroupRepo    *repo.ContactGroup
	ContactGroupService service.IContactGroupService
	ContactService      service.IContactService
}

func (g *Group) List(ctx context.Context, in *web.ContactGroupListRequest) (*web.ContactGroupListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	items := make([]*web.ContactGroupListResponse_Item, 0)

	count, err := g.ContactRepo.FindCount(ctx, "user_id = ? and status = ?", uid, model.Yes)
	if err != nil {
		return nil, err
	}

	items = append(items, &web.ContactGroupListResponse_Item{
		Name:  "全部",
		Count: int32(count),
	})

	group, err := g.ContactGroupService.GetUserGroup(ctx, uid)
	if err != nil {
		return nil, err
	}

	for _, v := range group {
		items = append(items, &web.ContactGroupListResponse_Item{
			Id:    int32(v.Id),
			Name:  v.Name,
			Count: int32(v.Num),
			Sort:  int32(v.Sort),
		})
	}

	return &web.ContactGroupListResponse{Items: items}, nil
}

func (g *Group) Save(ctx context.Context, in *web.ContactGroupSaveRequest) (*web.ContactGroupSaveResponse, error) {

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	updateItems := make([]*model.ContactGroup, 0)
	deleteItems := make([]int, 0)
	insertItems := make([]*model.ContactGroup, 0)

	ids := make(map[int]struct{})
	for i, item := range in.GetItems() {
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

	all, err := g.ContactGroupRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range all {
		if _, ok := ids[m.Id]; !ok {
			deleteItems = append(deleteItems, m.Id)
		}
	}

	if len(deleteItems) == 0 && len(updateItems) == 0 && len(insertItems) == 0 {
		return &web.ContactGroupSaveResponse{}, nil
	}

	err = g.ContactGroupRepo.Txx(ctx, func(tx *gorm.DB) error {

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
		return nil, err
	}

	return &web.ContactGroupSaveResponse{}, nil
}
