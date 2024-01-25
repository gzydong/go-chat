package group

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service/message"

	"go-chat/internal/service"
)

type Notice struct {
	GroupMemberRepo    *repo.GroupMember
	GroupNoticeRepo    *repo.GroupNotice
	GroupNoticeService service.IGroupNoticeService
	GroupMemberService service.IGroupMemberService
	Message            message.IService
	UsersRepo          *repo.Users
}

// CreateAndUpdate 添加或编辑群公告
func (c *Notice) CreateAndUpdate(ctx *core.Context) error {

	in := &web.GroupNoticeEditRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), int(in.GroupId), uid) {
		return ctx.ErrorBusiness("无权限操作")
	}

	var (
		msg string
		err error
	)

	if in.NoticeId == 0 {
		err = c.GroupNoticeService.Create(ctx.Ctx(), &service.GroupNoticeEditOpt{
			UserId:    uid,
			GroupId:   int(in.GroupId),
			NoticeId:  int(in.NoticeId),
			Title:     in.Title,
			Content:   in.Content,
			IsTop:     int(in.IsTop),
			IsConfirm: int(in.IsConfirm),
		})
		msg = "添加群公告成功！"
	} else {
		err = c.GroupNoticeService.Update(ctx.Ctx(), &service.GroupNoticeEditOpt{
			GroupId:   int(in.GroupId),
			NoticeId:  int(in.NoticeId),
			Title:     in.Title,
			Content:   in.Content,
			IsTop:     int(in.IsTop),
			IsConfirm: int(in.IsConfirm),
		})
		msg = "更新群公告成功！"
	}

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	userInfo, err := c.UsersRepo.FindByIdWithCache(ctx.Ctx(), uid)
	if err == nil {
		_ = c.Message.CreateGroupMessage(ctx.Ctx(), message.CreateGroupMessageOption{
			MsgType:  entity.ChatMsgTypeGroupNotice,
			FromId:   uid,
			ToFromId: int(in.GroupId),
			Extra: jsonutil.Encode(model.TalkRecordExtraGroupNotice{
				OwnerId:   uid,
				OwnerName: userInfo.Nickname,
				Title:     in.Title,
				Content:   in.Content,
			}),
		})
	}

	return ctx.Success(nil, msg)
}

// Delete 删除群公告
func (c *Notice) Delete(ctx *core.Context) error {

	in := &web.GroupNoticeDeleteRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.GroupNoticeService.Delete(ctx.Ctx(), int(in.GroupId), int(in.NoticeId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil, "群公告删除成功！")
}

// List 获取群公告列表(所有)
func (c *Notice) List(ctx *core.Context) error {

	in := &web.GroupNoticeListRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	// 判断是否是群成员
	if !c.GroupMemberRepo.IsMember(ctx.Ctx(), int(in.GroupId), ctx.UserId(), true) {
		return ctx.ErrorBusiness("无获取数据权限！")
	}

	all, err := c.GroupNoticeRepo.GetListAll(ctx.Ctx(), int(in.GroupId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	items := make([]*web.GroupNoticeListResponse_Item, 0)
	for i := 0; i < len(all); i++ {
		items = append(items, &web.GroupNoticeListResponse_Item{
			Id:           int32(all[i].Id),
			Title:        all[i].Title,
			Content:      all[i].Content,
			IsTop:        int32(all[i].IsTop),
			IsConfirm:    int32(all[i].IsConfirm),
			ConfirmUsers: all[i].ConfirmUsers,
			Avatar:       all[i].Avatar,
			CreatorId:    int32(all[i].CreatorId),
			CreatedAt:    timeutil.FormatDatetime(all[i].CreatedAt),
			UpdatedAt:    timeutil.FormatDatetime(all[i].UpdatedAt),
		})
	}

	return ctx.Success(&web.GroupNoticeListResponse{
		Items: items,
	})
}
