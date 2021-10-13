package helper

import (
	"github.com/gin-gonic/gin"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go-chat/config"
)

// Claims 相关信息
type Claims struct {
	Guard  string `json:"guard"`
	UserId int    `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJwtToken 创建 Jwt Token 的函数
// guard 登录守卫,可区分不同的登录token
// id 登录用户ID
func GenerateJwtToken(conf *config.Config, guard string, id int) (map[string]interface{}, error) {
	// 过期时间
	expiredAt := time.Now().Add(time.Second * time.Duration(conf.Jwt.ExpiresTime)).Unix()

	claims := Claims{
		Guard:  guard,
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer:    "go-chat",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenClaims.SignedString([]byte(conf.Jwt.Secret))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token":      token,
		"expired_at": expiredAt,
	}, nil
}

// ParseJwtToken 解析 Jwt Token 参数信息
func ParseJwtToken(secret string, token string) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return claims, err
}

// GetAuthToken 获取登录授权 token
func GetAuthToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	token = strings.TrimLeft(token, "Bearer")
	token = strings.TrimSpace(token)

	// Headers 中没有授权信息则读取 url 中的 token
	if len(token) == 0 {
		token = c.DefaultQuery("token", "")
	}

	if len(token) == 0 {
		token = c.DefaultPostForm("token", "")
	}

	return token
}
