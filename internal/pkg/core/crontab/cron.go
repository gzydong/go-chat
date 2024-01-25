package crontab

import (
	"context"
	"reflect"
)

type ICrontab interface {
	Name() string

	// Spec 配置定时任务规则
	Spec() string

	// Enable 是否启动
	Enable() bool

	// Do 任务执行入口
	Do(ctx context.Context) error
}

func ToCrontab(value any) []ICrontab {

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
