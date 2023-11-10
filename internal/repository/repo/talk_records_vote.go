package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecordsVote struct {
	ichat.Repo[model.TalkRecordsVote]
	cache *cache.Vote
}

type VoteStatistics struct {
	Count   int            `json:"count"`
	Options map[string]int `json:"options"`
}

func NewTalkRecordsVote(db *gorm.DB, cache *cache.Vote) *TalkRecordsVote {
	return &TalkRecordsVote{Repo: ichat.NewRepo[model.TalkRecordsVote](db), cache: cache}
}

func (t *TalkRecordsVote) GetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
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

func (t *TalkRecordsVote) SetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	uids := make([]int, 0)

	err := t.Repo.Db.WithContext(ctx).Table("talk_records_vote_answer").Where("vote_id = ?", vid).Pluck("user_id", &uids).Error

	if err != nil {
		return nil, err
	}

	_ = t.cache.SetVoteAnswerUser(ctx, vid, uids)

	return uids, nil
}

func (t *TalkRecordsVote) GetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	value, err := t.cache.GetVoteStatistics(ctx, vid)
	if err != nil {
		return t.SetVoteStatistics(ctx, vid)
	}

	statistic := &VoteStatistics{}

	_ = jsonutil.Decode(value, statistic)

	return statistic, nil
}

func (t *TalkRecordsVote) SetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	var (
		vote         model.TalkRecordsVote
		answerOption map[string]any
		options      = make([]string, 0)
	)

	tx := t.Repo.Db.WithContext(ctx)
	if err := tx.Table("talk_records_vote").First(&vote, vid).Error; err != nil {
		return nil, err
	}

	if err := jsonutil.Decode(vote.AnswerOption, &answerOption); err != nil {
		return nil, err
	}

	err := tx.Table("talk_records_vote_answer").Where("vote_id = ?", vid).Pluck("option", &options).Error
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
