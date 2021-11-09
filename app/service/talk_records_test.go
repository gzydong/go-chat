package service

import (
	"context"
	"fmt"
	"go-chat/app/pkg/jsonutil"
	"go-chat/testutil"
	"testing"
)

func testTalkRecordsService() *TalkRecordsService {
	return NewTalkRecordsService(NewBaseService(testutil.GetDb(), testutil.TestRedisClient()))
}

func TestTalkRecords(t *testing.T) {
	service := testTalkRecordsService()

	items, _ := service.GetTalkRecords(context.Background(), &QueryTalkRecordsOpts{
		TalkType:   1,
		UserId:     2054,
		ReceiverId: 2055,
		RecordId:   0,
		Limit:      30,
	})

	fmt.Println(jsonutil.JsonEncode(items))
}
