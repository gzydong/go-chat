package group

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/jsonutil"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
	"github.com/gzydong/go-chat/internal/service/message"
	"gorm.io/gorm"
)

var _ web.IGroupNoticeHandler = (*Notice)(nil)

type Notice struct {
	GroupMemberRepo    *repo.GroupMember
	GroupNoticeRepo    *repo.GroupNotice
	GroupMemberService service.IGroupMemberService
	Message            message.IService
	UsersRepo          *repo.Users
}

// Edit 添加或编辑群公告
func (c *Notice) Edit(ctx context.Context, in *web.GroupNoticeEditRequest) (*web.GroupNoticeEditResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if !c.GroupMemberRepo.IsMember(ctx, int(in.GroupId), uid, false) {
		return nil, entity.ErrPermissionDenied
	}

	notice, err := c.GroupNoticeRepo.FindByWhere(ctx, "group_id = ?", in.GroupId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if notice == nil {
		err = c.GroupNoticeRepo.Create(ctx, &model.GroupNotice{
			GroupId:      int(in.GroupId),
			CreatorId:    uid,
			ModifyId:     uid,
			Content:      in.Content,
			ConfirmUsers: "[]",
			IsConfirm:    2,
		})
	} else {
		_, err = c.GroupNoticeRepo.UpdateByWhere(ctx, map[string]any{
			"content":    in.Content,
			"modify_id":  uid,
			"updated_at": time.Now(),
		}, "group_id = ?", in.GroupId)
	}

	if err != nil {
		return nil, err
	}

	userInfo, err := c.UsersRepo.FindByIdWithCache(ctx, uid)
	if err == nil {
		_ = c.Message.CreateGroupMessage(ctx, message.CreateGroupMessageOption{
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

	return &web.GroupNoticeEditResponse{}, nil
}
