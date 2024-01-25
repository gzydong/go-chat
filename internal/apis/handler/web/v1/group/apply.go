package group

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/business"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"gorm.io/gorm"
)

type Apply struct {
	Redis              *redis.Client
	GroupApplyStorage  *cache.GroupApplyStorage
	GroupRepo          *repo.Group
	GroupApplyRepo     *repo.GroupApply
	GroupMemberRepo    *repo.GroupMember
	GroupApplyService  service.IGroupApplyService
	GroupMemberService service.IGroupMemberService
	GroupService       service.IGroupService
	PushMessage        *business.PushMessage
}

func (c *Apply) Create(ctx *core.Context) error {
	in := &web.GroupApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	apply, err := c.GroupApplyRepo.FindByWhere(ctx.Ctx(), "group_id = ? and user_id = ? and status = ?", in.GroupId, ctx.UserId(), model.GroupApplyStatusWait)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Error(err.Error())
	}

	uid := ctx.UserId()

	applyId := 0
	if apply == nil {
		data := &model.GroupApply{
			GroupId: int(in.GroupId),
			UserId:  uid,
			Status:  model.GroupApplyStatusWait,
			Remark:  in.Remark,
		}

		err = c.GroupApplyRepo.Create(ctx.Ctx(), data)
		if err == nil {
			applyId = data.Id
		}
	} else {
		applyId = apply.Id
		data := map[string]any{
			"remark":     in.Remark,
			"updated_at": timeutil.DateTime(),
		}

		_, err = c.GroupApplyRepo.UpdateByWhere(ctx.Ctx(), data, "id = ?", apply.Id)
	}

	if err != nil {
		return ctx.Error(err.Error())
	}

	find, err := c.GroupMemberRepo.FindByWhere(ctx.Ctx(), "group_id = ? and leader = ?", in.GroupId, model.GroupMemberLeaderOwner)
	if err == nil && find != nil {
		c.GroupApplyStorage.Incr(ctx.Ctx(), find.UserId)
	}

	_ = c.PushMessage.Push(ctx.Ctx(), entity.ImChannelChat, &entity.SubscribeMessage{
		Event: entity.SubEventGroupApply,
		Payload: jsonutil.Encode(entity.SubEventGroupApplyPayload{
			GroupId: int(in.GroupId),
			UserId:  ctx.UserId(),
			ApplyId: applyId,
		}),
	})

	return ctx.Success(nil)
}

func (c *Apply) Agree(ctx *core.Context) error {
	uid := ctx.UserId()

	in := &web.GroupApplyAgreeRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	apply, err := c.GroupApplyRepo.FindById(ctx.Ctx(), int(in.ApplyId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Error(err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.ErrorBusiness("申请信息不存在")
	}

	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), apply.GroupId, uid) {
		return ctx.Forbidden("无权限访问")
	}

	if apply.Status != model.GroupApplyStatusWait {
		return ctx.ErrorBusiness("申请信息已被他(她)人处理")
	}

	if !c.GroupMemberRepo.IsMember(ctx.Ctx(), apply.GroupId, apply.UserId, false) {
		err = c.GroupService.Invite(ctx.Ctx(), &service.GroupInviteOpt{
			UserId:    uid,
			GroupId:   apply.GroupId,
			MemberIds: []int{apply.UserId},
		})

		if err != nil {
			return ctx.ErrorBusiness(err.Error())
		}
	}

	data := map[string]any{
		"status":     model.GroupApplyStatusPass,
		"updated_at": timeutil.DateTime(),
	}

	_, err = c.GroupApplyRepo.UpdateByWhere(ctx.Ctx(), data, "id = ?", in.ApplyId)
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(nil)
}

func (c *Apply) Decline(ctx *core.Context) error {
	in := &web.GroupApplyDeclineRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	apply, err := c.GroupApplyRepo.FindById(ctx.Ctx(), int(in.ApplyId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Error(err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.ErrorBusiness("申请信息不存在")
	}

	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), apply.GroupId, uid) {
		return ctx.Forbidden("无权限访问")
	}

	if apply.Status != model.GroupApplyStatusWait {
		return ctx.ErrorBusiness("申请信息已被他(她)人处理")
	}

	data := map[string]any{
		"status":     model.GroupApplyStatusRefuse,
		"reason":     in.Remark,
		"updated_at": timeutil.DateTime(),
	}

	_, err = c.GroupApplyRepo.UpdateByWhere(ctx.Ctx(), data, "id = ?", in.ApplyId)
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(&web.GroupApplyDeclineResponse{})
}

func (c *Apply) List(ctx *core.Context) error {
	in := &web.GroupApplyListRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), int(in.GroupId), ctx.UserId()) {
		return ctx.Forbidden("无权限访问")
	}

	list, err := c.GroupApplyRepo.List(ctx.Ctx(), []int{int(in.GroupId)})
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

func (c *Apply) All(ctx *core.Context) error {
	uid := ctx.UserId()

	all, err := c.GroupMemberRepo.FindAll(ctx.Ctx(), func(db *gorm.DB) {
		db.Select("group_id")
		db.Where("user_id = ?", uid)
		db.Where("leader in ?", []int{
			model.GroupMemberLeaderOwner,
			model.GroupMemberLeaderAdmin,
		})
		db.Where("is_quit = ?", model.No)
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
		c.GroupApplyStorage.Del(ctx.Ctx(), ctx.UserId())
		return ctx.Success(resp)
	}

	list, err := c.GroupApplyRepo.List(ctx.Ctx(), groupIds)
	if err != nil {
		return ctx.ErrorBusiness("系统异常，请稍后再试！")
	}

	groups, err := c.GroupRepo.FindAll(ctx.Ctx(), func(db *gorm.DB) {
		db.Select("id,name")
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

	c.GroupApplyStorage.Del(ctx.Ctx(), ctx.UserId())

	return ctx.Success(resp)
}

func (c *Apply) ApplyUnreadNum(ctx *core.Context) error {
	return ctx.Success(map[string]any{
		"unread_num": c.GroupApplyStorage.Get(ctx.Ctx(), ctx.UserId()),
	})
}
