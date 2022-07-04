package cron

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	crontab "go-chat/internal/job/internal/handle/cron"
)

type Command *cli.Command

type ICrontab interface {
	// Spec 配置定时任务规则
	Spec() string

	// Handle 任务执行入口
	Handle(ctx context.Context) error
}

// Subcommands 注册的任务请务必实现 ICrontab 接口
type Subcommands struct {
	ClearWsCache      *crontab.ClearWsCache
	ClearArticle      *crontab.ClearArticle
	ClearTmpFile      *crontab.ClearTmpFile
	ClearExpireServer *crontab.ClearExpireServer
}

func NewCrontabCommand(handles *Subcommands) Command {
	return &cli.Command{
		Name:  "crontab",
		Usage: "Crontab Command - 常驻定时任务",
		Action: func(ctx *cli.Context) error {
			c := cron.New()

			for _, job := range toCrontab(handles) {
				_, err := c.AddFunc(job.Spec(), func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Printf("ICrontab err: %v \n", err)
						}
					}()

					_ = job.Handle(ctx.Context)
				})

				if err != nil {
					panic(err)
				}
			}

			log.Println("Crontab 定时任务已启动...")

			return run(c, ctx.Context)
		},
	}
}

func run(cron *cron.Cron, ctx context.Context) error {

	cron.Start()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	select {
	case <-s:
		cron.Stop()
	case <-ctx.Done():
		cron.Stop()
	}

	fmt.Println()
	log.Println("Crontab 定时任务已关闭")

	return nil
}

func toCrontab(value interface{}) []ICrontab {
	jobs := make([]ICrontab, 0)
	elem := reflect.ValueOf(value).Elem()
	for i := 0; i < elem.NumField(); i++ {
		if v, ok := elem.Field(i).Interface().(ICrontab); ok {
			jobs = append(jobs, v)
		}
	}

	return jobs
}
