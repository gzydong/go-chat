package group

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ginutil"
	"go-chat/internal/pkg/timeutil"

	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/service"
)

type Notice struct {
	service *service.GroupNoticeService
	member  *service.GroupMemberService
}

func NewGroupNoticeHandler(service *service.GroupNoticeService, member *service.GroupMemberService) *Notice {
	return &Notice{service: service, member: member}
}

// CreateAndUpdate 添加或编辑群公告
func (c *Notice) CreateAndUpdate(ctx *gin.Context) error {
	params := &web.GroupNoticeEditRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	var (
		msg string
		err error
	)

	uid := jwtutil.GetUid(ctx)

	if !c.member.Dao().IsLeader(params.GroupId, uid) {
		return ginutil.BusinessError(ctx, "无权限操作")
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
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, nil, msg)
}

// Delete 删除群公告
func (c *Notice) Delete(ctx *gin.Context) error {
	params := &web.GroupNoticeCommonRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	if err := c.service.Delete(ctx, params.GroupId, params.NoticeId); err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, nil, "群公告删除成功！")
}

// List 获取群公告列表(所有)
func (c *Notice) List(ctx *gin.Context) error {
	params := &web.GroupNoticeListRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	// 判断是否是群成员
	if !c.member.Dao().IsMember(params.GroupId, jwtutil.GetUid(ctx), true) {
		return ginutil.BusinessError(ctx, "无获取数据权限！")
	}

	items, _ := c.service.Dao().GetListAll(ctx, params.GroupId)

	rows := make([]*entity.H, 0)
	for i := 0; i < len(items); i++ {
		row := entity.H{}
		row["id"] = items[i].Id
		row["title"] = items[i].Title
		row["content"] = items[i].Content
		row["is_top"] = items[i].IsTop
		row["is_confirm"] = items[i].IsConfirm
		row["created_at"] = timeutil.FormatDatetime(items[i].CreatedAt)
		row["updated_at"] = timeutil.FormatDatetime(items[i].UpdatedAt)
		row["creator_id"] = items[i].CreatorId
		row["confirm_users"] = items[i].ConfirmUsers
		row["avatar"] = items[i].Avatar
		row["nickname"] = items[i].Nickname
		rows = append(rows, &row)
	}

	return ginutil.Success(ctx, entity.H{
		"rows": rows,
	})
}
