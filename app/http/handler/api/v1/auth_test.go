package v1

import (
	"fmt"
	"net/url"
	"testing"

	"go-chat/provider"

	"github.com/gin-gonic/gin"
	"go-chat/app/cache"

	"github.com/stretchr/testify/assert"
	"go-chat/app/dao"
	"go-chat/app/service"
	"go-chat/testutil"
)

func testAuth() *Auth {
	config := testutil.GetConfig()
	db := provider.NewMySQLClient(config)
	redisClient := testutil.TestRedisClient()
	user := &dao.UserDao{Base: &dao.Base{db: db}}

	userService := service.NewUserService(user)
	smsService := service.NewSmsService(&cache.SmsCodeCache{Redis: redisClient})
	authTokenCache := &cache.AuthTokenCache{Redis: redisClient}
	lockCache := cache.NewRedisLock(redisClient)

	return NewAuthHandler(
		config,
		userService,
		smsService,
		authTokenCache,
		lockCache,
	)
}

func TestAuth_Login(t *testing.T) {
	a := testAuth()
	r := testutil.NewTestRequest("/auth/login", a.Login)
	value := &url.Values{}
	value.Add("mobile", "18798272054")
	value.Add("password", "admin123")
	value.Add("platform", "windows")
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
		context.Set("__user_id__", 2054)
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
