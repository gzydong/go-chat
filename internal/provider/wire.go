package provider

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// 基础服务
	NewMySQLClient,
	NewRedisClient,
	NewHttpClient,
	NewEmailClient,
	NewFilesystem,
	NewBase64Captcha,
	NewIpAddressClient,
	NewRsa,
	wire.Struct(new(Providers), "*"),
)
