package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type GroupService struct {
	db *gorm.DB
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{
		db: db,
	}
}

// Create 创建群聊
func (s *GroupService) Create(ctx *gin.Context, request *request.GroupCreateRequest) error {
	var (
		err error
	)

	// 登录用户ID
	UserId := auth.GetAuthUserID(ctx)
	datetime := timeutil.DateTime()

	// 群成员用户ID
	MembersIds := strutil.ParseIds(request.MembersIds)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		GroupModel := &model.Group{
			CreatorId: UserId,
			GroupName: request.GroupName,
			Profile:   request.Profile,
			Avatar:    request.Avatar,
			MaxNum:    200,
			CreatedAt: datetime,
		}

		if err = tx.Create(GroupModel).Error; err != nil {
			return err
		}

		members := make([]*model.GroupMember, 0)
		for _, MemberId := range MembersIds {
			leader := 0
			if UserId == MemberId {
				leader = 2
			}

			members = append(members, &model.GroupMember{
				GroupId:   GroupModel.ID,
				UserId:    MemberId,
				Leader:    leader,
				CreatedAt: time.Now(),
			})
		}

		if err = tx.Model(&model.GroupMember{}).Create(members).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// Dismiss 解散群组(群主权限)
func (s *GroupService) Dismiss(GroupId int, UserId int) error {
	var (
		err error
	)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		queryModel := &model.Group{ID: GroupId, CreatorId: UserId}
		dismissedAt := sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		if err = tx.Model(queryModel).Updates(model.Group{IsDismiss: 1, DismissedAt: dismissedAt}).Error; err != nil {
			return err
		}

		err = s.db.Model(&model.GroupMember{}).Where("group_id = ?", GroupId).Unscoped().Updates(model.GroupMember{
			IsQuit:    1,
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}).Error

		if err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})

	return err
}

// Secede 退出群组(仅普通管理员及群成员)
func (s *GroupService) Secede(GroupId int, UserId int) error {
	var err error

	info := &model.GroupMember{}
	err = s.db.Model(model.GroupMember{}).Where("group_id = ? AND user_id = ?", GroupId, UserId).Unscoped().First(info).Error
	if err != nil {
		return err
	}

	if info.Leader == 2 {
		return errors.New("群主不能退出群组！")
	}

	if info.IsQuit == 1 {
		return nil
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		count := tx.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ?", GroupId, UserId).Unscoped().Updates(model.GroupMember{
			IsQuit:    1,
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}).RowsAffected

		if count == 0 {
			return nil
		}

		// todo 添加群消息

		return nil
	})

	return err
}

// UpdateMemberCard 修改群名片
func (s *GroupService) UpdateMemberCard(groupId int, userId int, remarks string) error {
	return s.db.Model(model.GroupMember{}).
		Where("group_id = ? and user_id = ?", groupId, userId).
		Unscoped().
		Update("user_card", remarks).Error
}

type Result struct {
	Id        int    `json:"id"`
	GroupName string `json:"group_name"`
	Avatar    string `json:"avatar"`
	Profile   string `json:"profile"`
	Leader    int    `json:"leader"`
	IsDisturb int    `json:"is_disturb"`
}

func (s *GroupService) UserGroupList(userId int) ([]*Result, error) {
	var err error
	items := make([]*Result, 0)

	res := s.db.Table("lar_group_member").
		Select("lar_group.id,lar_group.group_name,lar_group.avatar,lar_group.profile,lar_group_member.leader").
		Joins("left join lar_group on lar_group.id = lar_group_member.group_id").
		Where("lar_group_member.user_id = ? and lar_group_member.is_quit = ?", userId, 0).
		Unscoped().
		Scan(&items)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return items, nil
	}

	ids := make([]int, res.RowsAffected)
	for _, item := range items {
		ids = append(ids, item.Id)
	}

	var list []map[string]interface{}
	err = s.db.Table("lar_talk_list").
		Select("receiver_id,is_disturb").
		Where("talk_type = ? and receiver_id in ?", 2, ids).Find(&list).Error
	if err != nil {
		return nil, err
	}

	lists, err := slice.ToMap(list, "receiver_id")
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if data, ok := lists[int64(item.Id)]; ok {
			val, _ := data["is_disturb"]
			item.IsDisturb = int(reflect.ValueOf(val).Int())
		}
	}

	return items, nil
}
