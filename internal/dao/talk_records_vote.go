package dao

import (
	"context"

	"go-chat/internal/cache"
	"go-chat/internal/model"
	"go-chat/internal/pkg/jsonutil"
)

type TalkRecordsVoteDao struct {
	*BaseDao
	cache *cache.TalkVote
}

type VoteStatistics struct {
	Count   int            `json:"count"`
	Options map[string]int `json:"options"`
}

func NewTalkRecordsVoteDao(base *BaseDao, cache *cache.TalkVote) *TalkRecordsVoteDao {
	return &TalkRecordsVoteDao{BaseDao: base, cache: cache}
}

// GetVoteAnswerUser
func (dao *TalkRecordsVoteDao) GetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	// 读取缓存
	if uids, err := dao.cache.GetVoteAnswerUser(ctx, vid); err == nil {
		return uids, nil
	}

	uids, err := dao.SetVoteAnswerUser(ctx, vid)
	if err != nil {
		return nil, err
	}

	return uids, nil
}

func (dao *TalkRecordsVoteDao) SetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	uids := make([]int, 0)

	err := dao.Db().Table("talk_records_vote_answer").Where("vote_id = ?", vid).Pluck("user_id", &uids).Error

	if err != nil {
		return nil, err
	}

	_ = dao.cache.SetVoteAnswerUser(ctx, vid, uids)

	return uids, nil
}

func (dao *TalkRecordsVoteDao) GetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	value, err := dao.cache.GetVoteStatistics(ctx, vid)
	if err != nil {
		return dao.SetVoteStatistics(ctx, vid)
	}

	statistic := &VoteStatistics{}

	_ = jsonutil.Decode(value, statistic)

	return statistic, nil
}

func (dao *TalkRecordsVoteDao) SetVoteStatistics(ctx context.Context, vid int) (*VoteStatistics, error) {
	var (
		err          error
		vote         *model.TalkRecordsVote
		answerOption map[string]interface{}
		options      = make([]string, 0)
	)

	if err = dao.Db().Table("talk_records_vote").First(&vote, vid).Error; err != nil {
		return nil, err
	}

	_ = jsonutil.Decode(vote.AnswerOption, &answerOption)

	err = dao.Db().Table("talk_records_vote_answer").Where("vote_id = ?", vid).Pluck("option", &options).Error
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

	_ = dao.cache.SetVoteStatistics(ctx, vid, jsonutil.Encode(statistic))

	return statistic, nil
}
