package jwtutil

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
)

var (
	ErrorNoLogin = errors.New("请登录后操作! ")
)

type IStore interface {
	// IsBlackList 判断是否是黑名单
	IsBlackList(ctx context.Context, token string) bool
}

// Auth 授权中间件
func Auth(secret string, guard string, store IStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetJwtToken(c)

		claims, err := check(guard, secret, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, entity.H{"code": 401, "message": err.Error()})
			c.Abort()
			return
		}

		// 这里还需要验证 token 黑名单
		if store.IsBlackList(context.Background(), token) {
			c.JSON(http.StatusUnauthorized, entity.H{"code": 401, "message": "请登录再试！"})
			c.Abort()
			return
		}

		// 设置登录用户ID
		uid, _ := strconv.Atoi(claims.Id)

		c.Set(uuid, uid)

		// 记录 jwt 相关信息
		c.Set("jwt", map[string]string{
			"token":      token,
			"expires_at": strconv.Itoa(int(claims.ExpiresAt)),
		})

		c.Next()
	}
}

func check(guard string, secret string, token string) (*AuthClaims, error) {
	if token == "" {
		return nil, ErrorNoLogin
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 判断权限认证守卫是否一致
	if claims.Valid() != nil || claims.Guard != guard {
		return nil, ErrorNoLogin
	}

	return claims, nil
}
