package helper

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Jwt 秘钥
var secret = []byte("836a9c3fea9bba4e0a4d51bd902fbcc5")

// JWT 相关信息
type Claims struct {
	Guard  string `json:"guard"`
	UserID int    `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJwtToken 创建 Jwt Token 的函数
// @params guard 登录守卫,可区分不同的登录token
// @params id 登录用户ID
func GenerateJwtToken(guard string, id int) (map[string]interface{}, error) {
	// 过期时间
	expiredAt := time.Now().Add(12 * time.Hour).Unix()

	claims := Claims{
		Guard:  guard,
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer:    "go-chat",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenClaims.SignedString(secret)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return map[string]interface{}{
		"token":      token,
		"expired_at": expiredAt,
	}, nil
}

// ParseJwtToken 解析 Jwt Token 参数信息
func ParseJwtToken(token string) (*Claims, error) {
	cla := &Claims{}

	_, err := jwt.ParseWithClaims(token, cla, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	return cla, err
}
