package group

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
)

type Apply struct {
	applyServ  *service.GroupApplyService
	memberServ *service.GroupMemberService
	groupServ  *service.GroupService
}

func NewApply(applyServ *service.GroupApplyService, memberServ *service.GroupMemberService, groupServ *service.GroupService) *Apply {
	return &Apply{applyServ: applyServ, memberServ: memberServ, groupServ: groupServ}
}

func (c *Apply) Create(ctx *ichat.Context) error {

	params := &web.GroupApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.applyServ.Insert(ctx.Ctx(), int(params.GroupId), ctx.UserId(), params.Remark)
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！")
	}

	// TODO 这里需要推送给群主

	return ctx.Success(entity.H{})
}

func (c *Apply) Agree(ctx *ichat.Context) error {

	params := &web.GroupApplyAgreeRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	apply := &model.GroupApply{}
	if err := c.applyServ.Db().First(apply, params.ApplyId).Error; err != nil {
		return ctx.ErrorBusiness("数据不存在！")
	}

	if !c.memberServ.Dao().IsLeader(apply.GroupId, uid) {
		return ctx.Forbidden("无权限访问")
	}

	if !c.memberServ.Dao().IsMember(apply.GroupId, apply.UserId, false) {
		err := c.groupServ.InviteMembers(ctx.Ctx(), &service.InviteGroupMembersOpt{
			UserId:    uid,
			GroupId:   apply.GroupId,
			MemberIds: []int{apply.UserId},
		})
		if err != nil {
			return ctx.ErrorBusiness("处理失败！")
		}
	}

	err := c.applyServ.Db().Delete(model.GroupApply{}, "id = ?", apply.Id).Error
	if err != nil {
		logger.Error("数据删除失败 err", err.Error())
	}

	return ctx.Success(entity.H{})
}

func (c *Apply) Delete(ctx *ichat.Context) error {

	params := &web.GroupApplyDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.applyServ.Delete(ctx.Ctx(), int(params.ApplyId), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！")
	}

	return ctx.Success(entity.H{})
}

func (c *Apply) List(ctx *ichat.Context) error {

	params := &web.GroupApplyListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.memberServ.Dao().IsLeader(int(params.GroupId), ctx.UserId()) {
		return ctx.Forbidden("无权限访问")
	}

	list, err := c.applyServ.Dao().List(ctx.Ctx(), int(params.GroupId))
	if err != nil {
		logger.Error("[Apply List] 接口异常 err:", err.Error())
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
