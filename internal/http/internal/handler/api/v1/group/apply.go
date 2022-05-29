package group

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/model"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/timeutil"
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

func (c *Apply) Create(ctx *gin.Context) {
	params := &request.GroupApplyCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.applyServ.Insert(ctx.Request.Context(), params.GroupId, jwtutil.GetUid(ctx), params.Remark)
	if err != nil {
		response.BusinessError(ctx, "创建群聊失败，请稍后再试！")
		return
	}

	// TODO 这里需要推送给群主

	response.Success(ctx, entity.H{})
}

func (c *Apply) Agree(ctx *gin.Context) {
	params := &request.GroupApplyAgreeRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	apply := &model.GroupApply{}
	if err := c.applyServ.Db().First(apply, params.ApplyId).Error; err != nil {
		response.BusinessError(ctx, "数据不存在！")
		return
	}

	if !c.memberServ.Dao().IsLeader(apply.GroupId, uid) {
		response.Unauthorized(ctx, "无权限访问")
		return
	}

	if !c.memberServ.Dao().IsMember(apply.GroupId, apply.UserId, false) {
		err := c.groupServ.InviteMembers(ctx, &service.InviteGroupMembersOpts{
			UserId:    uid,
			GroupId:   apply.GroupId,
			MemberIds: []int{apply.UserId},
		})
		if err != nil {
			response.BusinessError(ctx, "处理失败！")
			return
		}
	}

	err := c.applyServ.Db().Delete(model.GroupApply{}, "id = ?", apply.Id).Error
	if err != nil {
		logger.Error("数据删除失败 err", err.Error())
	}

	response.Success(ctx, entity.H{})
}

func (c *Apply) Delete(ctx *gin.Context) {
	params := &request.GroupApplyDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.applyServ.Delete(ctx, params.ApplyId, jwtutil.GetUid(ctx))
	if err != nil {
		response.BusinessError(ctx, "创建群聊失败，请稍后再试！")
		return
	}

	response.Success(ctx, entity.H{})
}

func (c *Apply) List(ctx *gin.Context) {
	params := &request.GroupApplyListRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if !c.memberServ.Dao().IsLeader(params.GroupId, jwtutil.GetUid(ctx)) {
		response.Unauthorized(ctx, "无权限访问")
		return
	}

	list, err := c.applyServ.Dao().List(ctx.Request.Context(), params.GroupId)
	if err != nil {
		logger.Error("[Apply List] 接口异常 err:", err.Error())
		response.BusinessError(ctx, "创建群聊失败，请稍后再试！")
		return
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

	response.Success(ctx, items)
}
