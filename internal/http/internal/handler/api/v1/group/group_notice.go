package group

import (
	"github.com/gin-gonic/gin"

	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/service"
)

type GroupNotice struct {
	service *service.GroupNoticeService
	member  *service.GroupMemberService
}

func NewGroupNoticeHandler(service *service.GroupNoticeService, member *service.GroupMemberService) *GroupNotice {
	return &GroupNotice{service: service, member: member}
}

// CreateAndUpdate 添加或编辑群公告
func (c *GroupNotice) CreateAndUpdate(ctx *gin.Context) {
	params := &request.GroupNoticeEditRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	var (
		msg string
		err error
	)

	uid := jwtutil.GetUid(ctx)

	if !c.member.Dao().IsLeader(params.GroupId, uid) {
		response.BusinessError(ctx, "无权限操作")
		return
	}

	if params.NoticeId == 0 {
		err = c.service.Create(ctx.Request.Context(), &service.GroupNoticeEditOpts{
			UserId:    uid,
			GroupId:   params.GroupId,
			NoticeId:  params.NoticeId,
			Title:     params.Title,
			Content:   params.Content,
			IsTop:     params.IsTop,
			IsConfirm: params.IsConfirm,
		})
		msg = "添加群公告成功！"
	} else {
		err = c.service.Update(ctx.Request.Context(), &service.GroupNoticeEditOpts{
			GroupId:   params.GroupId,
			NoticeId:  params.NoticeId,
			Title:     params.Title,
			Content:   params.Content,
			IsTop:     params.IsTop,
			IsConfirm: params.IsConfirm,
		})
		msg = "更新群公告成功！"
	}

	if err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil, msg)
	}
}

// Delete 删除群公告
func (c *GroupNotice) Delete(ctx *gin.Context) {
	params := &request.GroupNoticeCommonRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Delete(ctx, params.GroupId, params.NoticeId); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil, "群公告删除成功！")
	}
}

// List 获取群公告列表(所有)
func (c *GroupNotice) List(ctx *gin.Context) {
	params := &request.GroupNoticeListRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 判断是否是群成员
	if !c.member.Dao().IsMember(params.GroupId, jwtutil.GetUid(ctx), true) {
		response.BusinessError(ctx, "无获取数据权限！")
		return
	}

	response.Success(ctx, gin.H{
		"rows": c.service.List(ctx, params.GroupId),
	})
}
