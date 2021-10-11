package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-chat/app/repository"
	"go-chat/connect"
	"go-chat/testutil"
	"net/url"
	"testing"
)

func testUser() *User {
	conf := testutil.GetConfig()

	return &User{UserRepo: &repository.UserRepository{
		DB: connect.MysqlConnect(conf),
	}}
}

func TestUser_Detail(t *testing.T) {
	a := testUser()

	r := testutil.NewTestRequest("/user/detail", func(context *gin.Context) {
		context.Set("__user_id__", 1)
	}, a.Detail)

	value := &url.Values{}

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestUser_ChangeDetail(t *testing.T) {
	a := testUser()

	r := testutil.NewTestRequest("/user/change/detail", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.ChangeDetail)

	value := &url.Values{}
	value.Add("nickname", "返税款1")
	value.Add("gender", "1")

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestUser_ChangePassword(t *testing.T) {
	a := testUser()

	r := testutil.NewTestRequest("/user/change/password", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.ChangePassword)

	value := &url.Values{}
	value.Add("old_password", "admin123")
	value.Add("new_password", "admin123")

	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}
