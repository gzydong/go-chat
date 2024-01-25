package mission

import (
	"context"
	"fmt"
	"go-chat/internal/mission/cron"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	ctb "github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/pkg/core/crontab"
)

type CronProvider struct {
	Config  *config.Config
	Crontab *cron.Crontab
}

func Cron(ctx *cli.Context, app *CronProvider) error {
	c := ctb.New()

	for _, v := range crontab.ToCrontab(app.Crontab) {
		job := v

		_, _ = c.AddFunc(job.Spec(), func() {
			defer func() {
				if err := recover(); err != nil {
					slog.Log(ctx.Context, slog.LevelError, fmt.Sprintf("panic crontab %s %s", job.Name(), job.Spec()))
				}
			}()

			_ = job.Do(ctx.Context)
		})

		slog.Log(ctx.Context, slog.LevelInfo, fmt.Sprintf("start crontab %s [%s]", job.Name(), job.Spec()))
	}

	return run(c, ctx.Context)
}

func run(cron *ctb.Cron, ctx context.Context) error {

	cron.Start()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	select {
	case <-s:
		cron.Stop()
	case <-ctx.Done():
		cron.Stop()
	}

	slog.Log(ctx, slog.LevelInfo, "Crontab 定时任务已关闭")

	return nil
}
