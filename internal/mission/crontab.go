package mission

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-chat/internal/mission/cron"
	"go-chat/internal/pkg/logger"

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
					slog.ErrorContext(ctx.Context, fmt.Sprintf("panic crontab %s", job.Name()))
					logger.Errorf(fmt.Sprintf("panic crontab %s error: %s", job.Name(), err.(error).Error()))
				}
			}()

			if err := job.Do(ctx.Context); err != nil {
				logger.Errorf(fmt.Sprintf("crontab %s %s error: %s", job.Name(), job.Spec(), err.Error()))
				slog.ErrorContext(ctx.Context, fmt.Sprintf("crontab %s %s error: %s", job.Name(), job.Spec(), err.Error()))
			}
		})

		slog.InfoContext(ctx.Context, fmt.Sprintf("add crontab %s [%s]", job.Name(), job.Spec()))
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
