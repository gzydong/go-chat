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

type TalkVote struct {
	redis *redis.Client
}

func NewTalkVote(rds *redis.Client) *TalkVote {
	return &TalkVote{redis: rds}
}

func (t *TalkVote) GetVoteAnswerUser(ctx context.Context, voteId int) ([]int, error) {
	val, err := t.redis.Get(ctx, fmt.Sprintf(VoteUsersCache, voteId)).Result()

	if err != nil {
		return nil, err
	}

	var uids []int
	if err := jsonutil.Decode(val, &uids); err != nil {
		return nil, err
	}

	return uids, nil
}

func (t *TalkVote) SetVoteAnswerUser(ctx context.Context, vid int, uids []int) error {
	return t.redis.Set(ctx, fmt.Sprintf(VoteUsersCache, vid), jsonutil.Encode(uids), time.Hour*24).Err()
}

func (t *TalkVote) GetVoteStatistics(ctx context.Context, vid int) (string, error) {
	return t.redis.Get(ctx, fmt.Sprintf(VoteStatisticCache, vid)).Result()
}

func (t *TalkVote) SetVoteStatistics(ctx context.Context, vid int, value string) error {
	return t.redis.Set(ctx, fmt.Sprintf(VoteStatisticCache, vid), value, time.Hour*24).Err()
}
