package v1

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"go-chat/app/cache"

	"github.com/stretchr/testify/assert"
	"go-chat/app/repository"
	"go-chat/app/service"
	"go-chat/connect"
	"go-chat/testutil"
)

func testAuth() *Auth {
	conf := testutil.GetConfig()
	db := connect.MysqlConnect(conf)
	redisClient := testutil.TestRedisClient()
	smsService := &service.SmsService{SmsCodeCache: &cache.SmsCodeCache{Redis: redisClient}}
	user := &repository.UserRepository{DB: db}
	s := &service.UserService{Repo: user}
	return &Auth{Conf: conf, UserService: s, SmsService: smsService}
}

func TestAuth_Login(t *testing.T) {
	a := testAuth()
	r := testutil.NewTestRequest("/auth/login", a.Login)
	value := &url.Values{}
	value.Add("username", "18953025089")
	value.Add("password", "admin123")
	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestAuth_Register(t *testing.T) {
	a := testAuth()
	r := testutil.NewTestRequest("/auth/register", a.Register)

	value := &url.Values{}
	value.Add("nickname", "测试账号昵称")
	value.Add("mobile", "18953025199")
	value.Add("password", "admin123")
	value.Add("sms_code", "000000")
	value.Add("platform", "mac")

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestAuth_Refresh(t *testing.T) {
	a := testAuth()
	r := testutil.NewTestRequest("/auth/refresh", func(context *gin.Context) {
		context.Set("__user_id__", 1)
	}, a.Refresh)

	value := &url.Values{}

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestAuth_Forget(t *testing.T) {
	a := testAuth()
	r := testutil.NewTestRequest("/auth/forget", func(context *gin.Context) {

	}, a.Forget)

	value := &url.Values{}
	value.Add("mobile", "18798272054")
	value.Add("password", "123456")
	value.Add("sms_code", "123456")

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}
