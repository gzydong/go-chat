package provider

import (
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jwtutil"
)

type UserJwtAuthorize jwtutil.IJwtAuthorize[entity.WebClaims]

func NewWebUserJwtAuthorize(config *config.Config) UserJwtAuthorize {
	return jwtutil.NewJwtAuthorize[entity.WebClaims]([]byte(config.Jwt.Secret))
}
