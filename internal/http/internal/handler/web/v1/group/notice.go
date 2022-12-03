package group

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"

	"go-chat/internal/service"
)

type Notice struct {
	service *service.GroupNoticeService
	member  *service.GroupMemberService
}

func NewNotice(service *service.GroupNoticeService, member *service.GroupMemberService) *Notice {
	return &Notice{service: service, member: member}
}

// CreateAndUpdate 添加或编辑群公告
func (c *Notice) CreateAndUpdate(ctx *ichat.Context) error {

	params := &web.GroupNoticeEditRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	var (
		msg string
		err error
	)

	uid := ctx.UserId()

	if !c.member.Dao().IsLeader(int(params.GroupId), uid) {
		return ctx.ErrorBusiness("无权限操作")
	}

	if params.NoticeId == 0 {
		err = c.service.Create(ctx.Ctx(), &service.GroupNoticeEditOpt{
			UserId:    uid,
			GroupId:   int(params.GroupId),
			NoticeId:  int(params.NoticeId),
			Title:     params.Title,
			Content:   params.Content,
			IsTop:     int(params.IsTop),
			IsConfirm: int(params.IsConfirm),
		})
		msg = "添加群公告成功！"
	} else {
		err = c.service.Update(ctx.Ctx(), &service.GroupNoticeEditOpt{
			GroupId:   int(params.GroupId),
			NoticeId:  int(params.NoticeId),
			Title:     params.Title,
			Content:   params.Content,
			IsTop:     int(params.IsTop),
			IsConfirm: int(params.IsConfirm),
		})
		msg = "更新群公告成功！"
	}

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil, msg)
}

// Delete 删除群公告
func (c *Notice) Delete(ctx *ichat.Context) error {

	params := &web.GroupNoticeDeleteRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.Delete(ctx.Ctx(), int(params.GroupId), int(params.NoticeId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil, "群公告删除成功！")
}

// List 获取群公告列表(所有)
func (c *Notice) List(ctx *ichat.Context) error {

	params := &web.GroupNoticeListRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// 判断是否是群成员
	if !c.member.Dao().IsMember(int(params.GroupId), ctx.UserId(), true) {
		return ctx.ErrorBusiness("无获取数据权限！")
	}

	items, _ := c.service.Dao().GetListAll(ctx.Ctx(), int(params.GroupId))

	rows := make([]*web.GroupNoticeListResponse_Item, 0)
	for i := 0; i < len(items); i++ {
		rows = append(rows, &web.GroupNoticeListResponse_Item{
			Id:           int32(items[i].Id),
			Title:        items[i].Title,
			Content:      items[i].Content,
			IsTop:        int32(items[i].IsTop),
			IsConfirm:    int32(items[i].IsConfirm),
			ConfirmUsers: items[i].ConfirmUsers,
			Avatar:       items[i].Avatar,
			CreatorId:    int32(items[i].CreatorId),
			CreatedAt:    timeutil.FormatDatetime(items[i].CreatedAt),
			UpdatedAt:    timeutil.FormatDatetime(items[i].UpdatedAt),
		})
	}

	return ctx.Success(&web.GroupNoticeListResponse{
		Items: rows,
	})
}
