package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/jsonutil"
)

const (
	VoteUsersCache     = "talk:vote:answer-users:%d"
	VoteStatisticCache = "talk:vote:statistic:%d"
)

type Vote struct {
	redis *redis.Client
}

func NewVote(rds *redis.Client) *Vote {
	return &Vote{redis: rds}
}

func (t *Vote) GetVoteAnswerUser(ctx context.Context, voteId int) ([]int, error) {
	val, err := t.redis.Get(ctx, fmt.Sprintf(VoteUsersCache, voteId)).Result()

	if err != nil {
		return nil, err
	}

	var ids []int
	if err := jsonutil.Decode(val, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func (t *Vote) SetVoteAnswerUser(ctx context.Context, vid int, uids []int) error {
	return t.redis.Set(ctx, fmt.Sprintf(VoteUsersCache, vid), jsonutil.Encode(uids), time.Hour*24).Err()
}

func (t *Vote) GetVoteStatistics(ctx context.Context, vid int) (string, error) {
	return t.redis.Get(ctx, fmt.Sprintf(VoteStatisticCache, vid)).Result()
}

func (t *Vote) SetVoteStatistics(ctx context.Context, vid int, value string) error {
	return t.redis.Set(ctx, fmt.Sprintf(VoteStatisticCache, vid), value, time.Hour*24).Err()
}
