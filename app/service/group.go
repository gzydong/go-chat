package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/slice"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type GroupService struct {
	db            *gorm.DB
	memberService *GroupMemberService
}

func NewGroupService(db *gorm.DB, memberService *GroupMemberService) *GroupService {
	return &GroupService{
		db:            db,
		memberService: memberService,
	}
}

func (s *GroupService) FindById(id int) (*model.Group, error) {
	info := &model.Group{}

	s.db.First(&info, id)

	return info, nil
}

// Create 创建群聊
func (s *GroupService) Create(ctx *gin.Context, request *request.GroupCreateRequest) error {
	var (
		err      error
		members  []*model.GroupMember
		talkList []*model.TalkList
	)

	// 登录用户ID
	uid := auth.GetAuthUserID(ctx)

	// 群成员用户ID
	MembersIds := slice.ParseIds(request.MembersIds)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		group := &model.Group{
			CreatorId: uid,
			GroupName: request.Name,
			Profile:   request.Profile,
			Avatar:    request.Avatar,
			MaxNum:    200,
			CreatedAt: time.Now(),
		}

		if err = tx.Create(group).Error; err != nil {
			return err
		}

		for _, val := range MembersIds {
			leader := 0
			if uid == val {
				leader = 2
			}

			members = append(members, &model.GroupMember{
				GroupId:   group.ID,
				UserId:    val,
				Leader:    leader,
				CreatedAt: time.Now(),
			})

			talkList = append(talkList, &model.TalkList{
				TalkType:   2,
				UserId:     val,
				ReceiverId: group.ID,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}

		if err = tx.Model(&model.GroupMember{}).Create(members).Error; err != nil {
			return err
		}

		if err = tx.Model(&model.TalkList{}).Create(talkList).Error; err != nil {
			return err
		}

		// 需要插入群邀请记录

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

		if err = s.db.Model(&model.GroupMember{}).Where("group_id = ?", GroupId).Unscoped().Updates(model.GroupMember{
			IsQuit:    1,
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}).Error; err != nil {
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

// InviteUsers 邀请用户加入群聊
func (s *GroupService) InviteUsers(groupId int, uid int, uids []int) error {
	var (
		err            error
		addMembers     []*model.GroupMember
		addTalkList    []*model.TalkList
		updateTalkList []int
		talkList       []*model.TalkList
	)

	m := make(map[int]int, 0)
	for _, value := range s.memberService.GetMemberIds(groupId) {
		m[value] = 1
	}

	listHash := make(map[int]*model.TalkList)
	s.db.Select("id", "user_id", "is_delete").Where("user_id in ? and receiver_id = ? and talk_type = 2", uids, groupId).Find(&talkList)
	for _, item := range talkList {
		listHash[item.UserId] = item
	}

	for _, value := range uids {
		if _, ok := m[value]; !ok {
			addMembers = append(addMembers, &model.GroupMember{
				GroupId:   groupId,
				UserId:    value,
				CreatedAt: time.Now(),
			})
		}

		if item, ok := listHash[value]; !ok {
			addTalkList = append(addTalkList, &model.TalkList{
				TalkType:   2,
				UserId:     value,
				ReceiverId: groupId,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		} else if item.IsDelete == 1 {
			updateTalkList = append(updateTalkList, item.ID)
		}
	}

	if len(addMembers) == 0 {
		return errors.New("邀请的好友，都已成为群成员")
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 删除已存在成员记录
		tx.Where("group_id = ? and user_id in ? and is_quit = 1", groupId, uids).Unscoped().Delete(model.GroupMember{})

		// 添加新成员
		if err = tx.Omit("deleted_at").Create(&addMembers).Error; err != nil {
			return err
		}

		// 添加用户的对话列表
		if len(addTalkList) > 0 {
			if err = tx.Select("talk_type", "user_id", "receiver_id", "updated_at").Create(&addTalkList).Error; err != nil {
				return err
			}
		}

		// 更新用户的对话列表
		if len(updateTalkList) > 0 {
			tx.Model(&model.TalkList{}).Where("id in ?", updateTalkList).Updates(map[string]interface{}{
				"is_delete":  0,
				"created_at": time.Now(),
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
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
