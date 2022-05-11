package jwtutil

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Options jwt.StandardClaims

const uuid = "__UID__"

type AuthClaims struct {
	Guard string `json:"guard"` // 授权守卫
	jwt.StandardClaims
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(guard string, secret string, ops *Options) string {
	claims := AuthClaims{
		Guard: guard,
		StandardClaims: jwt.StandardClaims{
			Audience:  ops.Audience,
			ExpiresAt: ops.ExpiresAt,
			Id:        ops.Id,
			IssuedAt:  ops.IssuedAt,
			Issuer:    ops.Issuer,
			NotBefore: ops.NotBefore,
			Subject:   ops.Subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString
}

// ParseToken 解析 JWT Token
func ParseToken(token string, secret string) (*AuthClaims, error) {
	claims := &AuthClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return claims, err
}

// GetUid 获取授权登录的用户ID
func GetUid(c *gin.Context) int {
	return c.GetInt(uuid)
}

// GetJwtToken 获取登录授权 token
func GetJwtToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))

	// Headers 中没有授权信息则读取 url 中的 token
	if token == "" {
		token = c.DefaultQuery("token", "")
	}

	if token == "" {
		token = c.DefaultPostForm("token", "")
	}

	return token
}
