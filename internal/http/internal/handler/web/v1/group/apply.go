package group

import (
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwtutil"
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

func NewApplyHandler(applyServ *service.GroupApplyService, memberServ *service.GroupMemberService, groupServ *service.GroupService) *Apply {
	return &Apply{applyServ: applyServ, memberServ: memberServ, groupServ: groupServ}
}

func (c *Apply) Create(ctx *ichat.Context) error {

	params := &web.GroupApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.applyServ.Insert(ctx.Context.Request.Context(), params.GroupId, jwtutil.GetUid(ctx.Context), params.Remark)
	if err != nil {
		return ctx.BusinessError("创建群聊失败，请稍后再试！")
	}

	// TODO 这里需要推送给群主

	return ctx.Success(entity.H{})
}

func (c *Apply) Agree(ctx *ichat.Context) error {

	params := &web.GroupApplyAgreeRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := jwtutil.GetUid(ctx.Context)

	apply := &model.GroupApply{}
	if err := c.applyServ.Db().First(apply, params.ApplyId).Error; err != nil {
		return ctx.BusinessError("数据不存在！")
	}

	if !c.memberServ.Dao().IsLeader(apply.GroupId, uid) {
		return ctx.Unauthorized("无权限访问")
	}

	if !c.memberServ.Dao().IsMember(apply.GroupId, apply.UserId, false) {
		err := c.groupServ.InviteMembers(ctx.Context, &service.InviteGroupMembersOpts{
			UserId:    uid,
			GroupId:   apply.GroupId,
			MemberIds: []int{apply.UserId},
		})
		if err != nil {
			return ctx.BusinessError("处理失败！")
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

	err := c.applyServ.Delete(ctx.Context, params.ApplyId, jwtutil.GetUid(ctx.Context))
	if err != nil {
		return ctx.WithMeta(map[string]interface{}{
			"error": err.Error(),
		}).BusinessError("创建群聊失败，请稍后再试！")
	}

	return ctx.Success(entity.H{})
}

func (c *Apply) List(ctx *ichat.Context) error {
	params := &web.GroupApplyListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.memberServ.Dao().IsLeader(params.GroupId, jwtutil.GetUid(ctx.Context)) {
		return ctx.Unauthorized("无权限访问")
	}

	list, err := c.applyServ.Dao().List(ctx.Context.Request.Context(), params.GroupId)
	if err != nil {
		logger.Error("[Apply List] 接口异常 err:", err.Error())
		return ctx.WithMeta(map[string]interface{}{
			"error": err.Error(),
		}).BusinessError("创建群聊失败，请稍后再试！")
	}

	items := make([]*entity.H, 0)
	for _, item := range list {
		items = append(items, &entity.H{
			"id":         item.Id,
			"user_id":    item.UserId,
			"group_id":   item.GroupId,
			"remark":     item.Remark,
			"avatar":     item.Avatar,
			"nickname":   item.Nickname,
			"created_at": timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	return ctx.Success(items)
}
