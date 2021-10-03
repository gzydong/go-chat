package v1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-chat/app/cache"
	"go-chat/app/repository"
	"go-chat/app/service"
	"go-chat/connect"
	"go-chat/testutil"
	"net/url"
	"testing"
)

func testCommon() *Common {
	conf := testutil.GetConfig()

	redisClient := testutil.TestRedisClient()
	smsService := &service.SmsService{SmsCodeCache: &cache.SmsCodeCache{Redis: redisClient}}

	return &Common{
		UserRepo:   &repository.UserRepository{DB: connect.MysqlConnect(conf)},
		SmsService: smsService,
	}
}

func TestCommon_SmsCode(t *testing.T) {
	common := testCommon()

	r := testutil.NewTestRequest("/common/sms-code", common.SmsCode)
	value := &url.Values{}
	value.Add("mobile", "18798276809")
	value.Add("channel", "login")

	resp, err := r.Form(value)
	assert.NoError(t, err)

	fmt.Println(resp.GetJson().Get("code"))
}
