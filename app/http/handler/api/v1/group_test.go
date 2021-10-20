package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-chat/app/service"
	"go-chat/provider"
	"go-chat/testutil"
	"net/url"
	"testing"
)

func testGroup() *Group {
	conf := testutil.GetConfig()
	db := provider.MysqlConnect(conf)
	groupService := service.NewGroupService(db)
	return NewGroupHandler(groupService)
}

func TestGroup_Create(t *testing.T) {
	a := testGroup()
	r := testutil.NewTestRequest("/api/v1/group/create", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.Create)

	value := &url.Values{}
	value.Add("group_name", "测试全发11")
	value.Add("ids", "1,2,3,4,5,6")
	value.Add("profile", "按手机卡那就看文")
	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestGroup_Dismiss(t *testing.T) {
	a := testGroup()
	r := testutil.NewTestRequest("/api/v1/group/dismiss", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.Dismiss)

	value := &url.Values{}
	value.Add("group_id", "280")
	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestGroup_Secede(t *testing.T) {
	a := testGroup()
	r := testutil.NewTestRequest("/api/v1/group/secede", func(context *gin.Context) {
		context.Set("__user_id__", 6)
	}, a.Secede)

	value := &url.Values{}
	value.Add("group_id", "280")
	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestGroup_EditGroupCard(t *testing.T) {
	a := testGroup()
	r := testutil.NewTestRequest("/api/v1/group/edit-card", func(context *gin.Context) {
		context.Set("__user_id__", 6)
	}, a.EditGroupCard)

	value := &url.Values{}
	value.Add("group_id", "280")
	value.Add("visit_card", "测试备注888")
	resp, err := r.Form(value)
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}

func TestGroup_GetGroups(t *testing.T) {
	a := testGroup()
	r := testutil.NewTestRequest("/api/v1/group/list", func(context *gin.Context) {
		context.Set("__user_id__", 2054)
	}, a.GetGroups)

	resp, err := r.Get()
	assert.NoError(t, err)
	fmt.Println(resp.GetJson().Get("code"))
}
