package job

import (
	"context"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/jedib0t/go-pretty/v6/table"
	crontab "github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/job/cron"
)

type ICrontab interface {
	Name() string

	// Spec 配置定时任务规则
	Spec() string

	// Enable 是否启动
	Enable() bool

	// Handle 任务执行入口
	Handle(ctx context.Context) error
}

type CronProvider struct {
	Config  *config.Config
	Crontab *Crontab
}

type Crontab struct {
	ClearWsCache      *cron.ClearWsCache
	ClearArticle      *cron.ClearArticle
	ClearTmpFile      *cron.ClearTmpFile
	ClearExpireServer *cron.ClearExpireServer
}

func Cron(ctx *cli.Context, app *CronProvider) error {
	c := crontab.New()

	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"#", "Name", "Time"})

	for i, exec := range toCrontab(app.Crontab) {
		job := exec

		_, _ = c.AddFunc(job.Spec(), func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Crontab Err: %v \n", err)
				}
			}()

			_ = job.Handle(ctx.Context)
		})

		tbl.AppendRow([]any{i + 1, job.Name(), job.Spec()})
	}

	tbl.Render()

	return run(c, ctx.Context)
}

func toCrontab(value any) []ICrontab {

	var jobs []ICrontab
	elem := reflect.ValueOf(value).Elem()
	for i := 0; i < elem.NumField(); i++ {
		if v, ok := elem.Field(i).Interface().(ICrontab); ok {
			if v.Enable() {
				jobs = append(jobs, v)
			}
		}
	}

	return jobs
}

func run(cron *crontab.Cron, ctx context.Context) error {

	cron.Start()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	select {
	case <-s:
		cron.Stop()
	case <-ctx.Done():
		cron.Stop()
	}

	log.Println("Crontab 定时任务已关闭")

	return nil
}
