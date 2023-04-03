package service

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
)

type GroupService struct {
	*repo.Source
	group    *repo.Group
	member   *repo.GroupMember
	relation *cache.Relation
	sequence *repo.Sequence
}

func NewGroupService(source *repo.Source, repo *repo.Group, member *repo.GroupMember, relation *cache.Relation, sequence *repo.Sequence) *GroupService {
	return &GroupService{Source: source, group: repo, member: member, relation: relation, sequence: sequence}
}

func (s *GroupService) Dao() *repo.Group {
	return s.group
}

type CreateGroupOpt struct {
	UserId    int    // 操作人ID
	Name      string // 群名称
	Avatar    string // 群头像
	Profile   string // 群简介
	MemberIds []int  // 联系人ID
}

// Create 创建群聊
func (s *GroupService) Create(ctx context.Context, opt *CreateGroupOpt) (int, error) {
	var (
		err      error
		members  []*model.GroupMember
		talkList []*model.TalkSession
	)

	// 群成员用户ID
	uids := sliceutil.Unique(append(opt.MemberIds, opt.UserId))

	group := &model.Group{
		CreatorId: opt.UserId,
		Name:      opt.Name,
		Profile:   opt.Profile,
		Avatar:    opt.Avatar,
		MaxNum:    model.GroupMemberMaxNum,
	}

	joinTime := time.Now()

	err = s.Db().Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(group).Error; err != nil {
			return err
		}

		for _, val := range uids {
			leader := 0
			if opt.UserId == val {
				leader = 2
			}

			members = append(members, &model.GroupMember{
				GroupId:  group.Id,
				UserId:   val,
				Leader:   leader,
				JoinTime: joinTime,
			})

			talkList = append(talkList, &model.TalkSession{
				TalkType:   2,
				UserId:     val,
				ReceiverId: group.Id,
			})
		}

		if err = tx.Create(members).Error; err != nil {
			return err
		}

		if err = tx.Create(talkList).Error; err != nil {
			return err
		}

		record := &model.TalkRecords{
			MsgId:      strutil.NewMsgId(),
			TalkType:   entity.ChatGroupMode,
			ReceiverId: group.Id,
			MsgType:    entity.MsgTypeGroupInvite,
			Sequence:   s.sequence.Get(ctx, 0, group.Id),
		}

		if err = tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})

	// 广播网关将在线的用户加入房间
	body := map[string]any{
		"event": entity.EventTalkJoinGroup,
		"data": jsonutil.Encode(map[string]any{
			"group_id": group.Id,
			"uids":     uids,
		}),
	}

	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))

	return group.Id, err
}

type UpdateGroupOpt struct {
	GroupId int    // 群ID
	Name    string // 群名称
	Avatar  string // 群头像
	Profile string // 群简介
}

// Update 更新群信息
func (s *GroupService) Update(ctx context.Context, opt *UpdateGroupOpt) error {

	_, err := s.group.UpdateById(ctx, opt.GroupId, map[string]any{
		"group_name": opt.Name,
		"avatar":     opt.Avatar,
		"profile":    opt.Profile,
	})

	return err
}

// Dismiss 解散群组[群主权限]
func (s *GroupService) Dismiss(ctx context.Context, groupId int, uid int) error {
	err := s.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Group{Id: groupId, CreatorId: uid}).Updates(&model.Group{
			IsDismiss: 1,
		}).Error; err != nil {
			return err
		}

		if err := s.Db().Model(&model.GroupMember{}).Where("group_id = ?", groupId).Updates(&model.GroupMember{
			IsQuit: 1,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// Secede 退出群组[仅管理员及群成员]
func (s *GroupService) Secede(ctx context.Context, groupId int, uid int) error {

	var info model.GroupMember
	if err := s.Db().Where("group_id = ? AND user_id = ? and is_quit = 0", groupId, uid).First(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("数据不存在！")
		}

		return err
	}

	if info.Leader == 2 {
		return errors.New("群主不能退出群组！")
	}

	var user model.Users
	err := s.Db().Table("users").Select("id,nickname").Where("id = ?", uid).First(&user).Error
	if err != nil {
		return err
	}

	record := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   entity.ChatGroupMode,
		ReceiverId: groupId,
		MsgType:    entity.MsgTypeGroupInvite,
		Sequence:   s.sequence.Get(ctx, 0, groupId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraGroupJoin{
			Action: 3,
			Operator: map[string]any{
				"uid":      user.Id,
				"nickname": user.Nickname,
			},
			Members: []map[string]any{},
		}),
	}

	err = s.Db().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ?", groupId, uid).Updates(&model.GroupMember{
			IsQuit: 1,
		}).Error
		if err != nil {
			return err
		}

		if err = tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.relation.DelGroupRelation(ctx, uid, groupId)

	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.EventTalkJoinGroup,
		"data": jsonutil.Encode(map[string]any{
			"type":     2,
			"group_id": groupId,
			"uids":     []int{uid},
		}),
	}))

	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.EventTalk,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"talk_type":   record.TalkType,
			"record_id":   record.Id,
		}),
	}))

	return nil
}

type InviteGroupMembersOpt struct {
	UserId    int   // 操作人ID
	GroupId   int   // 群ID
	MemberIds []int // 群成员ID
}

// InviteMembers 邀请加入群聊
func (s *GroupService) InviteMembers(ctx context.Context, opt *InviteGroupMembersOpt) error {
	var (
		err            error
		addMembers     []*model.GroupMember
		addTalkList    []*model.TalkSession
		updateTalkList []int
		talkList       []*model.TalkSession
	)

	m := make(map[int]struct{})
	for _, value := range s.member.GetMemberIds(ctx, opt.GroupId) {
		m[value] = struct{}{}
	}

	listHash := make(map[int]*model.TalkSession)
	s.Db().Select("id", "user_id", "is_delete").Where("user_id in ? and receiver_id = ? and talk_type = 2", opt.MemberIds, opt.GroupId).Find(&talkList)
	for _, item := range talkList {
		listHash[item.UserId] = item
	}

	mids := make([]int, 0)
	mids = append(mids, opt.MemberIds...)
	mids = append(mids, opt.UserId)

	memberItems := make([]*model.Users, 0)
	err = s.Db().Table("users").Select("id,nickname").Where("id in ?", mids).Scan(&memberItems).Error
	if err != nil {
		return err
	}

	memberMaps := make(map[int]*model.Users)
	for _, item := range memberItems {
		memberMaps[item.Id] = item
	}

	members := make([]map[string]any, 0)
	for _, value := range opt.MemberIds {
		member := memberMaps[value]
		members = append(members, map[string]any{
			"uid":      value,
			"nickname": member.Nickname,
		})

		if _, ok := m[value]; !ok {
			addMembers = append(addMembers, &model.GroupMember{
				GroupId:  opt.GroupId,
				UserId:   value,
				JoinTime: time.Now(),
			})
		}

		if item, ok := listHash[value]; !ok {
			addTalkList = append(addTalkList, &model.TalkSession{
				TalkType:   entity.ChatGroupMode,
				UserId:     value,
				ReceiverId: opt.GroupId,
			})
		} else if item.IsDelete == 1 {
			updateTalkList = append(updateTalkList, item.Id)
		}
	}

	if len(addMembers) == 0 {
		return errors.New("邀请的好友，都已成为群成员")
	}

	record := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   entity.ChatGroupMode,
		ReceiverId: opt.GroupId,
		MsgType:    entity.MsgTypeGroupInvite,
		Sequence:   s.sequence.Get(ctx, 0, opt.GroupId),
	}

	operator := memberMaps[opt.UserId]
	record.Extra = jsonutil.Encode(&model.TalkRecordExtraGroupJoin{
		Action: 1,
		Operator: map[string]any{
			"uid":      operator.Id,
			"nickname": operator.Nickname,
		},
		Members: members,
	})

	err = s.Db().Transaction(func(tx *gorm.DB) error {
		// 删除已存在成员记录
		tx.Where("group_id = ? and user_id in ? and is_quit = 1", opt.GroupId, opt.MemberIds).Delete(&model.GroupMember{})

		// 添加新成员
		if err = tx.Create(&addMembers).Error; err != nil {
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
			tx.Model(&model.TalkSession{}).Where("id in ?", updateTalkList).Updates(map[string]any{
				"is_delete":  0,
				"created_at": timeutil.DateTime(),
			})
		}

		if err = tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 广播网关将在线的用户加入房间
	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.EventTalkJoinGroup,
		"data": jsonutil.Encode(map[string]any{
			"type":     1,
			"group_id": opt.GroupId,
			"uids":     opt.MemberIds,
		}),
	}))

	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.EventTalk,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"talk_type":   record.TalkType,
			"record_id":   record.Id,
		}),
	}))

	return nil
}

type RemoveMembersOpt struct {
	UserId    int   // 操作人ID
	GroupId   int   // 群ID
	MemberIds []int // 群成员ID
}

// RemoveMembers 群成员移除群聊
func (s *GroupService) RemoveMembers(ctx context.Context, opt *RemoveMembersOpt) error {
	var num int64
	if err := s.Db().Model(&model.GroupMember{}).Where("group_id = ? and user_id in ? and is_quit = 0", opt.GroupId, opt.MemberIds).Count(&num).Error; err != nil {
		return err
	}

	if int(num) != len(opt.MemberIds) {
		return errors.New("删除失败")
	}

	mids := make([]int, 0)
	mids = append(mids, opt.MemberIds...)
	mids = append(mids, opt.UserId)

	memberItems := make([]*model.Users, 0)
	err := s.Db().Table("users").Select("id,nickname").Where("id in ?", mids).Scan(&memberItems).Error
	if err != nil {
		return err
	}

	memberMaps := make(map[int]*model.Users)
	for _, item := range memberItems {
		memberMaps[item.Id] = item
	}

	members := make([]map[string]any, 0)
	for _, value := range opt.MemberIds {
		member := memberMaps[value]
		members = append(members, map[string]any{
			"uid":      value,
			"nickname": member.Nickname,
		})
	}

	operator := memberMaps[opt.UserId]
	record := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		Sequence:   s.sequence.Get(ctx, 0, opt.GroupId),
		TalkType:   entity.ChatGroupMode,
		ReceiverId: opt.GroupId,
		MsgType:    entity.MsgTypeGroupInvite,
		Extra: jsonutil.Encode(&model.TalkRecordExtraGroupJoin{
			Action: 2,
			Operator: map[string]any{
				"uid":      operator.Id,
				"nickname": operator.Nickname,
			},
			Members: members,
		}),
	}

	err = s.Db().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.GroupMember{}).Where("group_id = ? and user_id in ? and is_quit = 0", opt.GroupId, opt.MemberIds).Updates(map[string]any{
			"is_quit":    1,
			"updated_at": time.Now(),
		}).Error
		if err != nil {
			return err
		}

		return tx.Create(record).Error
	})

	// 推送消息
	if err != nil {
		return err
	}

	s.relation.BatchDelGroupRelation(ctx, opt.MemberIds, opt.GroupId)

	_, _ = s.Redis().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
			"event": entity.EventTalkJoinGroup,
			"data": jsonutil.Encode(map[string]any{
				"type":     2,
				"group_id": opt.GroupId,
				"uids":     opt.MemberIds,
			}),
		}))

		pipe.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
			"event": entity.EventTalk,
			"data": jsonutil.Encode(map[string]any{
				"sender_id":   int64(record.UserId),
				"receiver_id": int64(record.ReceiverId),
				"talk_type":   record.TalkType,
				"record_id":   int64(record.Id),
			}),
		}))
		return nil
	})

	return nil
}

type session struct {
	ReceiverID int `json:"receiver_id"`
	IsDisturb  int `json:"is_disturb"`
}

func (s *GroupService) List(userId int) ([]*model.GroupItem, error) {
	tx := s.Db().Table("group_member")
	tx.Select("`group`.id,`group`.group_name,`group`.avatar,`group`.profile,group_member.leader")
	tx.Joins("left join `group` on `group`.id = group_member.group_id")
	tx.Where("group_member.user_id = ? and group_member.is_quit = ?", userId, 0)
	tx.Order("group_member.created_at desc")

	items := make([]*model.GroupItem, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	length := len(items)
	if length == 0 {
		return items, nil
	}

	ids := make([]int, 0, length)
	for i := range items {
		ids = append(ids, items[i].Id)
	}

	query := s.Db().Table("talk_session")
	query.Select("receiver_id,is_disturb")
	query.Where("talk_type = ? and receiver_id in ?", 2, ids)

	list := make([]*session, 0)
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	hash := make(map[int]*session)
	for i := range list {
		hash[list[i].ReceiverID] = list[i]
	}

	for i := range items {
		if value, ok := hash[items[i].Id]; ok {
			items[i].IsDisturb = value.IsDisturb
		}
	}

	return items, nil
}
