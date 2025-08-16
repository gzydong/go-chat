package cron

import "github.com/google/wire"

type Crontab struct {
	ClearArticle *ClearArticle
	ClearTmpFile *ClearTmpFile
}

var ProviderSet = wire.NewSet(
	wire.Struct(new(ClearArticle), "*"),
	wire.Struct(new(ClearTmpFile), "*"),
	wire.Struct(new(Crontab), "*"),
)
