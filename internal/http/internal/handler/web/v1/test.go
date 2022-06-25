package v1

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
)

type Test struct {
}

func NewTest() *Test {
	return &Test{}
}

func (c *Test) Success(ctx *ichat.Context) error {
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
