package v1

import (
	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/pkg/email"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/tmpl"
)

type Test struct {
	config      *config.Config
	emailClient *email.Client
}

func NewTest(config *config.Config, emailClient *email.Client) *Test {
	return &Test{config: config, emailClient: emailClient}
}

func (c *Test) Success(ctx *ichat.Context) error {

	fileContent, err := tmpl.Templates().ReadFile("resource/email/verify_code.tmpl")
	if err != nil {
		return err
	}

	body, _ := email.RenderString(string(fileContent), map[string]string{
		"code":         "123456",
		"service_name": "修改密码",
		"domain":       "https://im.gzydong.club",
	})

	_ = c.emailClient.SendMail(&email.Option{
		To:      []string{"837215079@qq.com"},
		Subject: "测试邮件",
		Body:    body,
	})

	return ctx.Success(&web.AuthLoginResponse{
		Type:        "",
		AccessToken: "",
		ExpiresIn:   15,
	})
}

func (c *Test) Raw(ctx *ichat.Context) error {
	return ctx.Raw("那框架是否那可就你那就开始DNA看就是那")
}

func (c *Test) Error(ctx *ichat.Context) error {
	return ctx.WithMeta(map[string]interface{}{
		"name": "maskjfank",
	}).BusinessError("业务错误")
}

func (c *Test) Invalid(ctx *ichat.Context) error {
	return ctx.InvalidParams("手机号不正确")
}

func (c *Test) WithData(ctx *ichat.Context) error {
	return ctx.WithMeta(map[string]interface{}{
		"name": "maskjfank",
	}).Error("maskjfank")
}
