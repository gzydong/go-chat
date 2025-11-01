package provider

import (
	"github.com/gzydong/go-chat/config"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/jwtutil"
)

type UserJwtAuthorize jwtutil.IJwtAuthorize[entity.WebClaims]

func NewWebUserJwtAuthorize(config *config.Config) UserJwtAuthorize {
	return jwtutil.NewJwtAuthorize[entity.WebClaims]([]byte(config.Jwt.Secret))
}
