package v1

import "go-chat/internal/pkg/ichat"

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

// Login 登录接口
func (c *Auth) Login(ctx *ichat.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *ichat.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *ichat.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}
