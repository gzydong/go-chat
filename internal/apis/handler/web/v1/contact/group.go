package contact

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"gorm.io/gorm"
)

type Group struct {
	ContactRepo         *repo.Contact
	ContactGroupRepo    *repo.ContactGroup
	ContactGroupService service.IContactGroupService
	ContactService      service.IContactService
}

// List 联系人分组列表
func (c *Group) List(ctx *core.Context) error {

	uid := ctx.AuthId()

	items := make([]*web.ContactGroupListResponse_Item, 0)

	count, err := c.ContactRepo.FindCount(ctx.GetContext(), "user_id = ? and status = ?", uid, model.Yes)
	if err != nil {
		return ctx.Error(err)
	}

	items = append(items, &web.ContactGroupListResponse_Item{
		Name:  "全部",
		Count: int32(count),
	})

	group, err := c.ContactGroupService.GetUserGroup(ctx.GetContext(), uid)
	if err != nil {
		return ctx.Error(err)
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

func (c *Group) Save(ctx *core.Context) error {
	in := &web.ContactGroupSaveRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

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

	all, err := c.ContactGroupRepo.FindAll(ctx.GetContext())
	if err != nil {
		return ctx.Error(err)
	}

	for _, m := range all {
		if _, ok := ids[m.Id]; !ok {
			deleteItems = append(deleteItems, m.Id)
		}
	}

	if len(deleteItems) == 0 && len(updateItems) == 0 && len(insertItems) == 0 {
		return ctx.Success(&web.ContactGroupSaveResponse{})
	}

	err = c.ContactGroupRepo.Txx(ctx.GetContext(), func(tx *gorm.DB) error {

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
		return ctx.Error(err)
	}

	return ctx.Success(&web.ContactGroupSaveResponse{})
}
