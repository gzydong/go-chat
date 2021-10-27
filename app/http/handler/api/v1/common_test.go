package v1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-chat/app/cache"
	"go-chat/app/repository"
	"go-chat/app/service"
	"go-chat/provider"
	"go-chat/testutil"
	"net/url"
	"testing"
)

func testCommon() *Common {
	conf := testutil.GetConfig()

	redisClient := testutil.TestRedisClient()
	smsService := service.NewSmsService(&cache.SmsCodeCache{Redis: redisClient})

	UserRepo := &dao.UserRepository{DB: provider.MysqlConnect(conf)}

	return NewCommonHandler(conf, smsService, UserRepo)
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

func TestCommon_EmailCode(t *testing.T) {
	//common := testCommon()
	//
	//r := testutil.NewTestRequest("/common/sms-code", common.SmsCode)
	//value := &url.Values{}
	//value.Add("mobile", "18798276809")
	//value.Add("channel", "login")
	//
	//resp, err := r.Form(value)
	//assert.NoError(t, err)
	//
	//fmt.Println(resp.GetJson().Get("code"))

	//em := &mail.EmailOptions{
	//	To:      []string{"837215079@qq.com"},
	//	Subject: "asjkfnaskjfa najksfna 测试挑剔madlkfmgad najskfnaek",
	//	Body:    "这是asmfa  没法看兰撒马卡龙问马拉喀什  麻辣可分为吗拉卡萨 马来开放门口啦 卡死了没法看论文吗拉菲马拉喀什吗发来看we马克里斯吗发案例 samdfajk  吗科利达麻辣阿斯达麻辣  卡里面分为 阿卡丽舒服吗马赛克发 马拉喀什吗发sklfa个测试项ergmslkdmlsdmfs",
	//}
	//
	//conf := testutil.GetConfig()
	//
	//fmt.Printf("%#v", conf.Email)
	//fmt.Println(em)
	//_ = mail.SendMail(conf.Email, em)
}
