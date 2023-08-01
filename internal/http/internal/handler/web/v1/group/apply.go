package group

import (
	"github.com/redis/go-redis/v9"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
	"gorm.io/gorm"
)

type Apply struct {
	groupApplyService  *service.GroupApplyService
	groupMemberService *service.GroupMemberService
	groupService       *service.GroupService
	storage            *cache.GroupApplyStorage
	redis              *redis.Client
}

func NewApply(groupApplyService *service.GroupApplyService, groupMemberService *service.GroupMemberService, groupService *service.GroupService, storage *cache.GroupApplyStorage, redis *redis.Client) *Apply {
	return &Apply{groupApplyService: groupApplyService, groupMemberService: groupMemberService, groupService: groupService, storage: storage, redis: redis}
}

func (c *Apply) Create(ctx *ichat.Context) error {

	params := &web.GroupApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.groupApplyService.Insert(ctx.Ctx(), int(params.GroupId), ctx.UserId(), params.Remark)
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！")
	}

	find, err := c.groupMemberService.Dao().FindByWhere(ctx.Ctx(), "group_id = ? and leader = ?", params.GroupId, 2)
	if err == nil && find != nil {
		c.storage.Incr(ctx.Ctx(), find.UserId)
	}

	c.redis.Publish(ctx.Ctx(), entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.SubEventGroupApply,
		"data": jsonutil.Encode(map[string]any{
			"group_id": params.GroupId,
			"user_id":  ctx.UserId(),
		}),
	}))

	return ctx.Success(nil)
}

func (c *Apply) Agree(ctx *ichat.Context) error {

	params := &web.GroupApplyAgreeRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	apply := &model.GroupApply{}
	if err := c.groupApplyService.Db().First(apply, params.ApplyId).Error; err != nil {
		return ctx.ErrorBusiness("数据不存在！")
	}

	if !c.groupMemberService.Dao().IsLeader(ctx.Ctx(), apply.GroupId, uid) {
		return ctx.Forbidden("无权限访问")
	}

	if !c.groupMemberService.Dao().IsMember(ctx.Ctx(), apply.GroupId, apply.UserId, false) {
		err := c.groupService.Invite(ctx.Ctx(), &service.GroupInviteOpt{
			UserId:    uid,
			GroupId:   apply.GroupId,
			MemberIds: []int{apply.UserId},
		})
		if err != nil {
			return ctx.ErrorBusiness("处理失败！")
		}
	}

	err := c.groupApplyService.Db().Delete(model.GroupApply{}, "id = ?", apply.Id).Error
	if err != nil {
		logger.Error("数据删除失败 err", err.Error())
	}

	return ctx.Success(nil)
}

func (c *Apply) Delete(ctx *ichat.Context) error {

	params := &web.GroupApplyDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.groupApplyService.Delete(ctx.Ctx(), int(params.ApplyId), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！" + err.Error())
	}

	return ctx.Success(nil)
}

func (c *Apply) List(ctx *ichat.Context) error {

	params := &web.GroupApplyListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.groupMemberService.Dao().IsLeader(ctx.Ctx(), int(params.GroupId), ctx.UserId()) {
		return ctx.Forbidden("无权限访问")
	}

	list, err := c.groupApplyService.Dao().List(ctx.Ctx(), []int{int(params.GroupId)})
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！")
	}

	items := make([]*web.GroupApplyListResponse_Item, 0)
	for _, item := range list {
		items = append(items, &web.GroupApplyListResponse_Item{
			Id:        int32(item.Id),
			UserId:    int32(item.UserId),
			GroupId:   int32(item.GroupId),
			Remark:    item.Remark,
			Avatar:    item.Avatar,
			Nickname:  item.Nickname,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	return ctx.Success(&web.GroupApplyListResponse{Items: items})
}

func (c *Apply) All(ctx *ichat.Context) error {

	uid := ctx.UserId()

	all, err := c.groupMemberService.Dao().FindAll(ctx.Ctx(), func(db *gorm.DB) {
		db.Select("group_id")
		db.Where("user_id = ?", uid)
		db.Where("leader = ?", 2)
		db.Where("is_quit = ?", 0)
	})

	if err != nil {
		return ctx.ErrorBusiness("系统异常，请稍后再试！")
	}

	groupIds := make([]int, 0, len(all))
	for _, m := range all {
		groupIds = append(groupIds, m.GroupId)
	}

	resp := &web.GroupApplyAllResponse{Items: make([]*web.GroupApplyAllResponse_Item, 0)}

	if len(groupIds) == 0 {
		return ctx.Success(resp)
	}

	list, err := c.groupApplyService.Dao().List(ctx.Ctx(), groupIds)
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！")
	}

	groups, err := c.groupService.Dao().FindAll(ctx.Ctx(), func(db *gorm.DB) {
		db.Select("id,group_name")
		db.Where("id in ?", groupIds)
	})
	if err != nil {
		return err
	}

	groupMap := sliceutil.ToMap(groups, func(t *model.Group) int {
		return t.Id
	})

	for _, item := range list {
		resp.Items = append(resp.Items, &web.GroupApplyAllResponse_Item{
			Id:        int32(item.Id),
			UserId:    int32(item.UserId),
			GroupName: groupMap[item.GroupId].Name,
			GroupId:   int32(item.GroupId),
			Remark:    item.Remark,
			Avatar:    item.Avatar,
			Nickname:  item.Nickname,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	c.storage.Del(ctx.Ctx(), ctx.UserId())

	return ctx.Success(resp)
}

// ApplyUnreadNum 获取群申请未读数
func (c *Apply) ApplyUnreadNum(ctx *ichat.Context) error {
	return ctx.Success(map[string]any{
		"unread_num": c.storage.Get(ctx.Ctx(), ctx.UserId()),
	})
}
