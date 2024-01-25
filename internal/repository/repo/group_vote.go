package repo

import (
	"context"

	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type GroupVote struct {
	core.Repo[model.GroupVote]
	cache *cache.Vote
}

type VoteStatistics struct {
	Count   int            `json:"count"`
	Options map[string]int `json:"options"`
}

func NewGroupVote(db *gorm.DB, cache *cache.Vote) *GroupVote {
	return &GroupVote{Repo: core.NewRepo[model.GroupVote](db), cache: cache}
}

func (t *GroupVote) GetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	// 读取缓存
	if uids, err := t.cache.GetVoteAnswerUser(ctx, vid); err == nil {
		return uids, nil
	}

	uids, err := t.SetVoteAnswerUser(ctx, vid)
	if err != nil {
		return nil, err
	}

	return uids, nil
}

func (t *GroupVote) SetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	uids := make([]int, 0)

	err := t.Repo.Db.WithContext(ctx).Table("group_vote_answer").Where("vote_id = ?", vid).Pluck("user_id", &uids).Error

	if err != nil {
		return nil, err
	}

	_ = t.cache.SetVoteAnswerUser(ctx, vid, uids)

	return uids, nil
}

func (t *GroupVote) GetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	value, err := t.cache.GetVoteStatistics(ctx, vid)
	if err != nil {
		return t.SetVoteStatistics(ctx, vid)
	}

	statistic := &VoteStatistics{}

	_ = jsonutil.Decode(value, statistic)

	return statistic, nil
}

func (t *GroupVote) SetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	var (
		vote         model.GroupVote
		answerOption map[string]any
		options      = make([]string, 0)
	)

	tx := t.Repo.Db.WithContext(ctx)
	if err := tx.Table("group_vote").First(&vote, vid).Error; err != nil {
		return nil, err
	}

	if err := jsonutil.Decode(vote.AnswerOption, &answerOption); err != nil {
		return nil, err
	}

	err := tx.Table("group_vote_answer").Where("vote_id = ?", vid).Pluck("option", &options).Error
	if err != nil {
		return nil, err
	}

	opts := make(map[string]int)
	for option := range answerOption {
		opts[option] = 0
	}

	for _, option := range options {
		opts[option] += 1
	}

	statistic := &VoteStatistics{
		Options: opts,
		Count:   len(options),
	}

	_ = t.cache.SetVoteStatistics(ctx, vid, jsonutil.Encode(statistic))

	return statistic, nil
}

func (t *GroupVote) FindAllAnsweredUserList(ctx context.Context, viteId int) ([]model.GroupVoteAnswer, error) {
	var items []model.GroupVoteAnswer

	err := t.Repo.Db.Table("group_vote_answer").Where("vote_id = ?", viteId).Scan(&items).Error
	if err == nil {
		return items, nil
	}

	return nil, err
}
