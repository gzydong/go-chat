package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Options jwt.RegisteredClaims

type AuthClaims struct {
	Guard string `json:"guard"` // 授权守卫
	jwt.RegisteredClaims
}

func NewNumericDate(t time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(t)
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(guard string, secret string, ops *Options) string {

	claims := AuthClaims{
		Guard: guard,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  ops.Audience,
			ExpiresAt: ops.ExpiresAt,
			ID:        ops.ID,
			IssuedAt:  ops.IssuedAt,
			Issuer:    ops.Issuer,
			NotBefore: ops.NotBefore,
			Subject:   ops.Subject,
		},
	}

	tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	return tokenString
}

// ParseToken 解析 JWT Token
func ParseToken(token string, secret string) (*AuthClaims, error) {

	data, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if claims, ok := data.Claims.(*AuthClaims); ok && data.Valid {
		return claims, nil
	}

	return nil, err
}
