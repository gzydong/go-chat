package jwtutil

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims[T any] struct {
	Metadata T `json:"metadata"`
	jwt.RegisteredClaims
}

func NewTokenWithClaims[T any](secret []byte, metadata T, opts ...func(c *jwt.RegisteredClaims)) (string, error) {
	claims := JwtClaims[T]{
		Metadata: metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        strings.ReplaceAll(uuid.NewString(), "-", ""),
		},
	}

	for _, opt := range opts {
		opt(&claims.RegisteredClaims)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func WithTokenExpiresAt(expire time.Duration) func(c *jwt.RegisteredClaims) {
	return func(c *jwt.RegisteredClaims) {
		c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expire))
	}
}

func ParseWithClaims[T any](secret []byte, tokenString string) (*JwtClaims[T], error) {
	var claims = new(JwtClaims[T])
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*JwtClaims[T]); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}
