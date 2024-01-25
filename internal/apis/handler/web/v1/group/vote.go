package group

import (
	"time"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
)

type Vote struct {
	GroupMemberRepo  *repo.GroupMember
	GroupVoteRepo    *repo.GroupVote
	GroupVoteService service.IGroupVoteService
	MessageService   message.IService
}

// Create 创建投票
func (v *Vote) Create(ctx *core.Context) error {
	var in web.GroupVoteCreateRequest
	if err := ctx.Context.ShouldBind(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if len(in.Options) <= 1 {
		return ctx.InvalidParams("options 选项必须大于1！")
	}

	if len(in.Options) > 6 {
		return ctx.InvalidParams("options 选项不能超过6个！")
	}

	isAnonymous := false
	if in.IsAnonymous == 1 {
		isAnonymous = true
	}

	voteId, err := v.GroupVoteService.Create(ctx.Context, &service.GroupVoteCreateOpt{
		UserId:        uid,
		Title:         in.Title,
		AnswerMode:    int(in.Mode),
		AnswerOptions: in.Options,
		IsAnonymous:   isAnonymous,
		GroupId:       int(in.GroupId),
	})
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	if err := v.MessageService.CreateVoteMessage(ctx.Ctx(), message.CreateVoteMessage{
		TalkMode: entity.ChatGroupMode,
		FromId:   uid,
		ToFromId: int(in.GroupId),
		VoteId:   voteId,
	}); err != nil {
		logger.Errorf("创建投票消息失败：%v", err)
	}

	return ctx.Success(web.GroupVoteCreateResponse{})
}

// Submit 提交投票
func (v *Vote) Submit(ctx *core.Context) error {
	var in web.GroupVoteSubmitRequest
	if err := ctx.Context.ShouldBind(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := v.GroupVoteService.Submit(ctx.Context, &service.GroupVoteSubmitOpt{
		UserId:  ctx.UserId(),
		VoteId:  int(in.VoteId),
		Options: in.Options,
	})

	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(web.GroupVoteSubmitResponse{})
}

// Detail 投票详情
func (v *Vote) Detail(ctx *core.Context) error {
	var in web.GroupVoteDetailRequest
	if err := ctx.Context.ShouldBind(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	voteInfo, err := v.GroupVoteRepo.FindById(ctx.Ctx(), int(in.VoteId))
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	if !v.GroupMemberRepo.IsMember(ctx.Ctx(), voteInfo.GroupId, ctx.UserId(), false) {
		return ctx.Forbidden("暂无查看投票详情权限！")
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
	if err := jsonutil.Decode(voteInfo.AnswerOption, &options); err != nil {
		return ctx.Error(err.Error())
	}

	for _, option := range options {
		resp.AnswerOptions = append(resp.AnswerOptions, &web.GroupVoteDetailResponse_AnswerOption{
			Value: option.Value,
			Key:   option.Key,
		})
	}

	items, err := v.GroupVoteRepo.FindAllAnsweredUserList(ctx.Ctx(), voteInfo.Id)
	if err != nil {
		return ctx.Error(err.Error())
	}

	userId := ctx.UserId()
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

		for uid, item := range hashMap {
			resp.AnsweredUsers = append(resp.AnsweredUsers, item)

			if uid == userId {
				resp.IsSubmit = true
			}
		}
	}

	return ctx.Success(resp)
}
