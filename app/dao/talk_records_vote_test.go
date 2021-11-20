package dao

import (
	"context"
	"fmt"
	"go-chat/app/cache"
	"go-chat/testutil"
	"testing"
)

func TestTalkRecordsVoteDao(t *testing.T) {

	db := testutil.GetDb()
	rds := testutil.TestRedisClient()

	dao := NewTalkRecordsVoteDao(NewBaseDao(db), cache.NewTalkVote(rds))

	uids, _ := dao.GetVoteAnswerUser(context.Background(), 34)

	fmt.Println(uids)
}
