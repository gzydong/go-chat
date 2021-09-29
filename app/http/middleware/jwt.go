package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"go-chat/app/http/response"
	"go-chat/config"
)

// JwtAuth 授权中间件
func JwtAuth(conf *config.Config, guard string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helper.GetAuthToken(c)

		info, err := jwtCheckLogin(conf, guard, token)
		if err != nil {
			response.Unauthorized(c, err)
			c.Abort()
			return
		}

		// todo 黑名单判断
		// ...

		// 设置登录用户ID
		c.Set("__user_id__", info.UserId)

		c.Next()
	}
}

// jwtCheckLogin 验证登录
func jwtCheckLogin(conf *config.Config, guard string, token string) (*helper.Claims, error) {
	if token == "" {
		return nil, errors.New("请登录后操作! ")
	}

	jwt, err := helper.ParseJwtToken(conf, token)
	if err != nil {
		return nil, errors.New("Token 信息验证失败! ")
	}

	// 判断权限认证守卫是否一致
	if jwt.Guard != guard {
		return nil, errors.New("非法操作! ")
	}

	return jwt, nil
}
