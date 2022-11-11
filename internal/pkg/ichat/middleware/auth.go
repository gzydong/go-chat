package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/jwt"
)

const JWTSessionConst = "__JWT_SESSION__"

var (
	ErrorNoLogin = errors.New("请登录后操作! ")
)

type IStore interface {
	// IsBlackList 判断是否是黑名单
	IsBlackList(ctx context.Context, token string) bool
}

type JSession struct {
	Uid       int    `json:"uid"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// Auth 授权中间件
func Auth(secret string, guard string, store IStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := QueryToken(c)

		claims, err := verify(guard, secret, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
			c.Abort()
			return
		}

		// 这里还需要验证 token 黑名单
		if store.IsBlackList(c.Request.Context(), token) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请登录再试."})
			c.Abort()
			return
		}

		// 设置登录用户ID
		uid, err := strconv.Atoi(claims.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解析 jwt 失败."})
			c.Abort()
			return
		}

		// 记录 jwt 相关信息
		c.Set(JWTSessionConst, &JSession{
			Uid:       uid,
			Token:     token,
			ExpiresAt: claims.ExpiresAt.Unix(),
		})

		c.Next()
	}
}

func QueryToken(c *gin.Context) string {

	token := c.GetHeader("Authorization")
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))

	// Headers 中没有授权信息则读取 url 中的 token
	if token == "" {
		token = c.DefaultQuery("token", "")
	}

	return token
}

func verify(guard string, secret string, token string) (*jwt.AuthClaims, error) {

	if token == "" {
		return nil, ErrorNoLogin
	}

	claims, err := jwt.ParseToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 判断权限认证守卫是否一致
	if claims.Guard != guard || claims.Valid() != nil {
		return nil, ErrorNoLogin
	}

	return claims, nil
}
