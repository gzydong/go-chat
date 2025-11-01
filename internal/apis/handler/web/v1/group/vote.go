package group

import (
	"context"
	"time"

	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/errorx"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/jsonutil"
	"github.com/gzydong/go-chat/internal/pkg/logger"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
	"github.com/gzydong/go-chat/internal/service/message"
)

var _ web.IGroupVoteHandler = (*Vote)(nil)

type Vote struct {
	GroupMemberRepo  *repo.GroupMember
	GroupVoteRepo    *repo.GroupVote
	GroupVoteService service.IGroupVoteService
	MessageService   message.IService
}

// Create 创建投票
func (v *Vote) Create(ctx context.Context, in *web.GroupVoteCreateRequest) (*web.GroupVoteCreateResponse, error) {

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if len(in.Options) <= 1 {
		return nil, errorx.NewInvalidParams("options 选项必须大于1")
	}

	if len(in.Options) > 6 {
		return nil, errorx.NewInvalidParams("options 选项不能超过6个")
	}

	isAnonymous := in.IsAnonymous == 1

	voteId, err := v.GroupVoteService.Create(ctx, &service.GroupVoteCreateOpt{
		UserId:        uid,
		Title:         in.Title,
		AnswerMode:    int(in.Mode),
		AnswerOptions: in.Options,
		IsAnonymous:   isAnonymous,
		GroupId:       int(in.GroupId),
	})
	if err != nil {
		return nil, err
	}

	if err := v.MessageService.CreateVoteMessage(ctx, message.CreateVoteMessage{
		TalkMode: entity.ChatGroupMode,
		FromId:   uid,
		ToFromId: int(in.GroupId),
		VoteId:   voteId,
	}); err != nil {
		logger.Errorf("创建投票消息失败：%v", err)
	}

	return &web.GroupVoteCreateResponse{}, nil
}

// Submit 提交投票
func (v *Vote) Submit(ctx context.Context, in *web.GroupVoteSubmitRequest) (*web.GroupVoteSubmitResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	err := v.GroupVoteService.Submit(ctx, &service.GroupVoteSubmitOpt{
		UserId:  uid,
		VoteId:  int(in.VoteId),
		Options: in.Options,
	})

	if err != nil {
		return nil, err
	}

	return &web.GroupVoteSubmitResponse{}, nil
}

// Detail 投票详情
func (v *Vote) Detail(ctx context.Context, in *web.GroupVoteDetailRequest) (*web.GroupVoteDetailResponse, error) {
	voteInfo, err := v.GroupVoteRepo.FindById(ctx, int(in.VoteId))
	if err != nil {
		return nil, err
	}

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if !v.GroupMemberRepo.IsMember(ctx, voteInfo.GroupId, uid, false) {
		return nil, entity.ErrPermissionDenied
	}

	resp := &web.GroupVoteDetailResponse{
		VoteId:        int32(voteInfo.Id),
		Title:         voteInfo.Title,
		AnswerMode:    int32(voteInfo.AnswerMode),
		AnswerOptions: make([]*web.GroupVoteDetailResponse_AnswerOption, 0),
		AnswerNum:     int32(voteInfo.AnswerNum),
		AnsweredNum:   int32(voteInfo.AnsweredNum),
		IsAnonymous:   int32(voteInfo.IsAnonymous),
		AnsweredUsers: make([]*web.GroupVoteDetailResponse_AnsweredUser, 0),
		IsSubmit:      false,
	}

	var options []model.GroupVoteOption
	if err := jsonutil.Unmarshal(voteInfo.AnswerOption, &options); err != nil {
		return nil, err
	}

	for _, option := range options {
		resp.AnswerOptions = append(resp.AnswerOptions, &web.GroupVoteDetailResponse_AnswerOption{
			Value: option.Value,
			Key:   option.Key,
		})
	}

	items, err := v.GroupVoteRepo.FindAllAnsweredUserList(ctx, voteInfo.Id)
	if err != nil {
		return nil, err
	}

	userId := uid
	if len(items) > 0 {
		hashMap := make(map[int]*web.GroupVoteDetailResponse_AnsweredUser)

		for _, item := range items {
			if val, ok := hashMap[item.UserId]; !ok {
				hashMap[item.UserId] = &web.GroupVoteDetailResponse_AnsweredUser{
					UserId:     int32(item.UserId),
					Nickname:   "xxx",
					Options:    []string{item.Option},
					AnswerTime: item.CreatedAt.Format(time.DateTime),
				}
			} else {
				val.Options = append(val.Options, item.Option)
			}
		}

		for id, item := range hashMap {
			resp.AnsweredUsers = append(resp.AnsweredUsers, item)

			if id == userId {
				resp.IsSubmit = true
			}
		}
	}

	return resp, nil
}
