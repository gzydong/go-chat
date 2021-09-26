package v1

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/app/repository"
	"go-chat/app/service"
	"go-chat/connect"
	"go-chat/testutil"
)

func testAuth() *Auth {
	conf := testutil.GetConfig()
	db := connect.MysqlConnect(conf)
	user := &repository.UserRepository{DB: db}
	s := &service.UserService{Repo: user}
	return &Auth{Conf: conf, UserService: s}
}

func TestAuth_Login(t *testing.T) {
	a := testAuth()
	r := testutil.NewTestRequest("/auth/login", a.Login)
	value := &url.Values{}
	value.Add("username", "")
	value.Add("password", "")
	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetBody())
}
