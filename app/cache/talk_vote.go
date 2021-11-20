package cache

import "github.com/go-redis/redis/v8"

type TalkVote struct {
	rds *redis.Client
}

func NewTalkVote(rds *redis.Client) *TalkVote {
	return &TalkVote{rds: rds}
}

func (t *TalkVote) name() {

}
