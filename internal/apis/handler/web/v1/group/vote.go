package group

import (
	"time"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Vote struct {
	GroupMemberRepo  *repo.GroupMember
	GroupVoteRepo    *repo.GroupVote
	GroupVoteService service.IGroupVoteService
}

// Create 创建投票
func (v *Vote) Create(ctx *ichat.Context) error {
	var params web.GroupVoteCreateRequest
	if err := ctx.Context.ShouldBind(&params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if len(params.Options) <= 1 {
		return ctx.InvalidParams("options 选项必须大于1！")
	}

	if len(params.Options) > 6 {
		return ctx.InvalidParams("options 选项不能超过6个！")
	}

	isAnonymous := false
	if params.IsAnonymous == 1 {
		isAnonymous = true
	}

	if err := v.GroupVoteService.Create(ctx.Context, &service.GroupVoteCreateOpt{
		UserID:        uid,
		Title:         params.Title,
		AnswerMode:    int(params.Mode),
		AnswerOptions: params.Options,
		IsAnonymous:   isAnonymous,
		GroupId:       int(params.GroupId),
	}); err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(web.GroupVoteCreateResponse{})
}

// Submit 提交投票
func (v *Vote) Submit(ctx *ichat.Context) error {
	var params web.GroupVoteSubmitRequest
	if err := ctx.Context.ShouldBind(&params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := v.GroupVoteService.Submit(ctx.Context, &service.GroupVoteSubmitOpt{
		UserId:  ctx.UserId(),
		VoteId:  int(params.VoteId),
		Options: params.Options,
	})

	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(web.GroupVoteSubmitResponse{})
}

// Detail 投票详情
func (v *Vote) Detail(ctx *ichat.Context) error {
	var params web.GroupVoteSubmitRequest
	if err := ctx.Context.ShouldBind(&params); err != nil {
		return ctx.InvalidParams(err)
	}

	voteInfo, err := v.GroupVoteRepo.FindById(ctx.Ctx(), int(params.VoteId))
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
		AnswerNum:     int32(voteInfo.AnsweredNum),
		AnsweredNum:   int32(voteInfo.AnsweredNum),
		IsAnonymous:   int32(voteInfo.IsAnonymous),
		AnsweredUsers: make([]*web.GroupVoteDetailResponse_AnsweredUser, 0),
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

	if len(items) > 0 {
		hashMap := make(map[int]*web.GroupVoteDetailResponse_AnsweredUser)

		for _, item := range items {
			if val, ok := hashMap[item.UserId]; ok {
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

		for _, item := range hashMap {
			resp.AnsweredUsers = append(resp.AnsweredUsers, item)
		}
	}

	return ctx.Success(resp)
}
