package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/config"
	"strconv"
)

// JwtAuth 授权中间件
func JwtAuth(conf *config.Config, guard string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := auth.GetJwtToken(c)

		info, err := checkLogin(guard, conf.Jwt.Secret, token)
		if err != nil {
			response.Unauthorized(c, err)
			c.Abort()
			return
		}

		// 设置登录用户ID
		uid, _ := strconv.Atoi(info.Id)
		c.Set(entity.LoginUserID, uid)

		c.Next()
	}
}

func checkLogin(guard string, secret string, token string) (*auth.JwtAuthClaims, error) {
	if token == "" {
		return nil, errors.New("请登录后操作! ")
	}

	claims, err := auth.VerifyJwtToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 判断权限认证守卫是否一致
	if claims.Guard != guard {
		return nil, errors.New("非法操作! ")
	}

	return claims, nil
}
