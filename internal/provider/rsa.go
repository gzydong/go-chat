package provider

import (
	"go-chat/config"
	"go-chat/internal/pkg/encrypt/aesutil"
	"go-chat/internal/pkg/encrypt/rsautil"
)

func NewRsa(config *config.Config) rsautil.IRsa {
	return rsautil.NewRsa([]byte(config.App.PublicKey), []byte(config.App.PrivateKey))
}

func NewAesUtil(config *config.Config) aesutil.IAesUtil {
	if config.App.AesKey == "" {
		panic("app.aes_key is empty")
	}

	return aesutil.NewAesUtil(config.App.AesKey)
}
