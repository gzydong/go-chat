package service

import (
	"context"
	"errors"
	"fmt"

	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ IGroupVoteService = (*GroupVoteService)(nil)

type IGroupVoteService interface {
	Create(ctx context.Context, opt *GroupVoteCreateOpt) (int, error)
	Submit(ctx context.Context, opt *GroupVoteSubmitOpt) error
	Detail(ctx context.Context, opt *GroupVoteDetailOpt) error
}

type GroupVoteService struct {
	*repo.Source
	GroupMemberRepo *repo.GroupMember
	GroupVoteRepo   *repo.GroupVote
	Sequence        *repo.Sequence
}

type GroupVoteCreateOpt struct {
	GroupId       int      // 群组ID
	UserId        int      // 用户ID(创建人)
	Title         string   // 投票标题
	AnswerMode    int      // 答题模式[1:单选;2:多选;]
	AnswerOptions []string // 答题选项
	IsAnonymous   bool     // 匿名投票
}

func (g *GroupVoteService) Create(ctx context.Context, opt *GroupVoteCreateOpt) (int, error) {
	options := make([]model.GroupVoteOption, 0)
	for i, value := range opt.AnswerOptions {
		options = append(options, model.GroupVoteOption{
			Key:   fmt.Sprintf("%c", 65+i),
			Value: value,
		})
	}

	vote := &model.GroupVote{
		GroupId:      opt.GroupId,
		UserId:       opt.UserId,
		Title:        opt.Title,
		AnswerMode:   opt.AnswerMode,
		AnswerOption: jsonutil.Encode(options),
		AnswerNum:    int(g.GroupMemberRepo.CountMemberTotal(ctx, opt.GroupId)),
		Status:       model.VoteStatusWait,
	}

	if opt.IsAnonymous {
		vote.IsAnonymous = 1
	}

	if err := g.Source.Db().Create(vote).Error; err != nil {
		return 0, err
	}

	return vote.Id, nil
}

type GroupVoteSubmitOpt struct {
	UserId  int      // 用户ID(投票人)
	VoteId  int      // 投票ID
	Options []string // 投票选项
}

func (g *GroupVoteService) Submit(ctx context.Context, opt *GroupVoteSubmitOpt) error {
	db := g.Source.Db().WithContext(ctx)

	voteInfo, err := g.GroupVoteRepo.FindById(ctx, opt.VoteId)
	if err != nil {
		return err
	}

	if !g.GroupMemberRepo.IsMember(ctx, voteInfo.GroupId, opt.UserId, false) {
		return errors.New("暂无投票权限！")
	}

	var count int64
	db.Table("group_vote_answer").Where("vote_id = ? and user_id = ?", opt.VoteId, opt.UserId).Count(&count)
	if count > 0 {
		return fmt.Errorf("重复投票[%d]", opt.VoteId)
	}

	ops := opt.Options
	if voteInfo.AnswerMode == model.VoteAnswerModeSingle {
		ops = ops[:1]
	}

	err = g.Source.Db().Transaction(func(tx *gorm.DB) error {
		data := map[string]any{
			"answered_num": gorm.Expr("answered_num + 1"),
			"status":       gorm.Expr("if(answered_num >= answer_num, 1, 0)"),
		}

		if err := tx.Table("group_vote").Where("id = ?", voteInfo.Id).Updates(data).Error; err != nil {
			return err
		}

		answers := make([]*model.GroupVoteAnswer, 0, len(ops))
		for _, option := range ops {
			answers = append(answers, &model.GroupVoteAnswer{
				VoteId: voteInfo.Id,
				UserId: opt.UserId,
				Option: option,
			})
		}

		return tx.Create(answers).Error
	})

	if err != nil {
		return err
	}

	// TODO 投递消息

	return nil
}

type GroupVoteDetailOpt struct {
	UserId int // 用户ID
	VoteId int // 投票ID
}

func (g *GroupVoteService) Detail(ctx context.Context, opt *GroupVoteDetailOpt) error {
	return nil
}
