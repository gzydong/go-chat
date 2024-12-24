package provider

import (
	"go-chat/config"
	"go-chat/internal/pkg/encrypt/rsautil"
)

func NewRsa(config *config.Config) rsautil.IRsa {
	return rsautil.NewRsa([]byte(config.App.PublicKey), []byte(config.App.PrivateKey))
}
