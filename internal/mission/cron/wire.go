package cron

import "github.com/google/wire"

type Crontab struct {
	ClearWsCache      *ClearWsCache
	ClearArticle      *ClearArticle
	ClearTmpFile      *ClearTmpFile
	ClearExpireServer *ClearExpireServer
}

var ProviderSet = wire.NewSet(
	wire.Struct(new(ClearArticle), "*"),
	wire.Struct(new(ClearTmpFile), "*"),
	wire.Struct(new(ClearWsCache), "*"),
	wire.Struct(new(ClearExpireServer), "*"),
	wire.Struct(new(Crontab), "*"),
)
