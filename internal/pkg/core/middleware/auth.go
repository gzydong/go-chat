package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-chat/internal/pkg/jwtutil"
)

const JWTAuthID = "__JWT_AUTH_ID__"

type IStorage interface {
	// IsBlackList 判断是否是黑名单
	IsBlackList(ctx context.Context, token string) bool
}

type IClaims interface {
	GetAuthID() int
}

type AuthClaimsKey struct{}

func GetAuthToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	// Headers 中没有授权信息则读取 url 中的 token
	if token == "" {
		token = c.DefaultQuery("token", "")
	}

	return token
}

func NewJwtMiddleware[T IClaims](secret []byte, storage IStorage, fn func(ctx context.Context, claims *jwtutil.JwtClaims[T]) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetAuthToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "授权异常，请登录后操作!"})
			return
		}

		claims, err := jwtutil.ParseWithClaims[T](secret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
			return
		}

		if storage.IsBlackList(c.Request.Context(), token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "授权异常，请登录后操作!"})
			return
		}

		if err = fn(c.Request.Context(), claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
			return
		}

		// 将用户信息放入上下文
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), AuthClaimsKey{}, claims.Metadata))
		c.Set(JWTAuthID, claims.Metadata.GetAuthID())

		expiresAt := claims.ExpiresAt.Unix()
		// 提前15分钟刷新token
		if expiresAt-time.Now().Unix() < 900 {
			newToken, err := jwtutil.NewTokenWithClaims(secret, claims.Metadata, func(c *jwt.RegisteredClaims) {
				c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 2))
				c.Issuer = claims.Issuer
				c.Audience = claims.Audience
				c.Subject = claims.Subject
				c.NotBefore = claims.NotBefore
				c.IssuedAt = claims.IssuedAt
			})

			if err == nil {
				c.Header("Refresh-Access-Token", newToken)
				c.Header("Refresh-Access-Expires-At", "7200")
			}
		}

		c.Next()
	}
}

// FormContext 从上下文中获取用户信息
func FormContext[T IClaims](ctx context.Context) (T, error) {
	if ctx.Value(AuthClaimsKey{}) == nil {
		return *new(T), errors.New("claims is nil")
	}

	claims, ok := ctx.Value(AuthClaimsKey{}).(T)
	if !ok {
		return *new(T), errors.New("claims is nil")
	}

	return claims, nil
}
