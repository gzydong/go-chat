package group

import (
	"errors"
	"fmt"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
	"gorm.io/gorm"
	"time"
)

type Notice struct {
	GroupMemberRepo    *repo.GroupMember
	GroupNoticeRepo    *repo.GroupNotice
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

	if !c.GroupMemberRepo.IsMember(ctx.Ctx(), int(in.GroupId), uid, false) {
		return ctx.ErrorBusiness("无权限操作")
	}

	var (
		msg string
		err error
	)

	notice, err := c.GroupNoticeRepo.FindByWhere(ctx.Ctx(), "group_id = ?", in.GroupId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Error(err.Error())
	}

	if notice == nil {
		msg = "群公告创建成功！"
		err = c.GroupNoticeRepo.Create(ctx.Ctx(), &model.GroupNotice{
			GroupId:      int(in.GroupId),
			CreatorId:    ctx.UserId(),
			ModifyId:     ctx.UserId(),
			Content:      in.Content,
			ConfirmUsers: "[]",
			IsConfirm:    2,
		})
	} else {
		msg = "群公告更新成功！"
		_, err = c.GroupNoticeRepo.UpdateByWhere(ctx.Ctx(), map[string]any{
			"content":    in.Content,
			"modify_id":  ctx.UserId(),
			"updated_at": time.Now(),
		}, "group_id = ?", in.GroupId)
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
				Title:     fmt.Sprintf("【%s】 更新了群公告", userInfo.Nickname),
				Content:   in.Content,
			}),
		})
	}

	return ctx.Success(nil, msg)
}
