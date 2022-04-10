package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go-chat/internal/pkg/jsonutil"
)

const (
	VoteUsersCache     = "talk:vote:answer-users:%d"
	VoteStatisticCache = "talk:vote:statistic:%d"
)

type TalkVote struct {
	rds *redis.Client
}

func NewTalkVote(rds *redis.Client) *TalkVote {
	return &TalkVote{rds: rds}
}

func (cache *TalkVote) GetVoteAnswerUser(ctx context.Context, voteId int) ([]int, error) {
	val, err := cache.rds.Get(ctx, fmt.Sprintf(VoteUsersCache, voteId)).Result()

	if err != nil {
		return nil, err
	}

	var uids []int
	if err := jsonutil.Decode(val, &uids); err != nil {
		return nil, err
	}

	return uids, nil
}

func (cache *TalkVote) SetVoteAnswerUser(ctx context.Context, vid int, uids []int) error {
	return cache.rds.Set(ctx, fmt.Sprintf(VoteUsersCache, vid), jsonutil.Encode(uids), time.Hour*24).Err()
}

func (cache *TalkVote) GetVoteStatistics(ctx context.Context, vid int) (string, error) {
	return cache.rds.Get(ctx, fmt.Sprintf(VoteStatisticCache, vid)).Result()
}

func (cache *TalkVote) SetVoteStatistics(ctx context.Context, vid int, value string) error {
	return cache.rds.Set(ctx, fmt.Sprintf(VoteStatisticCache, vid), value, time.Hour*24).Err()
}
