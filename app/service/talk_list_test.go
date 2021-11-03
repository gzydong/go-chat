package service

import (
	"context"
	"go-chat/app/dao"
	"go-chat/testutil"
	"testing"
)

func TestTalkListService(t *testing.T) {
	db := testutil.GetDb()
	rds := testutil.TestRedisClient()

	ser := NewTalkListService(&BaseService{
		db:  db,
		rds: rds,
	}, dao.NewTalkListDao(db))

	ser.GetTalkList(context.Background())
}
