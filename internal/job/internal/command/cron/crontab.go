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

type CrontabHandle interface {
	// Spec 配置定时任务规则
	Spec() string

	// Handle 任务执行入口
	Handle(ctx context.Context) error
}

// Handles 注册的任务请务必实现 CrontabHandle 接口
type Handles struct {
	ClearWsCacheHandle      *crontab.ClearWsCacheHandle
	ClearArticleHandle      *crontab.ClearArticleHandle
	ClearTmpFileHandle      *crontab.ClearTmpFileHandle
	ClearExpireServerHandle *crontab.ClearExpireServerHandle
}

func NewCrontabCommand(handles *Handles) Command {
	return &cli.Command{
		Name:  "crontab",
		Usage: "Crontab Command | 常驻定时任务",
		Action: func(ctx *cli.Context) error {
			c := cron.New()

			for _, job := range toCrontabHandle(handles) {
				_, _ = c.AddFunc(job.Spec(), func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Printf("CrontabHandle err: %v \n", err)
						}
					}()

					_ = job.Handle(ctx.Context)
				})
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

func toCrontabHandle(value interface{}) []CrontabHandle {
	jobs := make([]CrontabHandle, 0)
	elem := reflect.ValueOf(value).Elem()
	for i := 0; i < elem.NumField(); i++ {
		if v, ok := elem.Field(i).Interface().(CrontabHandle); ok {
			jobs = append(jobs, v)
		}
	}

	return jobs
}
