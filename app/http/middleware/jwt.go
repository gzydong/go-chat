package middleware

import (
	"context"
	"errors"
	"go-chat/app/cache"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/config"
)

// JwtAuth 授权中间件
func JwtAuth(conf *config.Config, guard string, session *cache.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := auth.GetJwtToken(c)

		claims, err := check(guard, conf.Jwt.Secret, token)
		if err != nil {
			response.Unauthorized(c, err)
			c.Abort()
			return
		}

		// 这里还需要验证 token 黑名单
		if session.IsExistBlackList(context.Background(), token) {
			response.Unauthorized(c, "请登录再试！")
			c.Abort()
			return
		}

		// 设置登录用户ID
		uid, _ := strconv.Atoi(claims.Id)
		c.Set(entity.LoginUserID, uid)

		// 记录 jwt 相关信息
		c.Set("jwt", map[string]string{
			"token":      token,
			"expires_at": strconv.Itoa(int(claims.ExpiresAt)),
		})

		c.Next()
	}
}

func check(guard string, secret string, token string) (*auth.JwtAuthClaims, error) {
	if token == "" {
		return nil, errors.New("请登录后操作! ")
	}

	claims, err := auth.VerifyJwtToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 判断权限认证守卫是否一致
	if claims.Valid() != nil || claims.Guard != guard {
		return nil, errors.New("请登录后操作! ")
	}

	return claims, nil
}
