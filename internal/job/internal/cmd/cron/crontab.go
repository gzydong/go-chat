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
	"go-chat/internal/job/internal/handle/crontab"
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
		Usage: "定时任务",
		Action: func(ctx *cli.Context) error {
			c := cron.New()

			jobs := toCrontabHandle(handles)
			for _, job := range jobs {
				_, _ = c.AddFunc(job.Spec(), func() {
					_ = job.Handle(ctx.Context)
				})
			}

			fmt.Println("Crontab 定时任务已启动...")
			fmt.Println("当选任务数：", len(jobs))

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
