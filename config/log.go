package config

import (
	"fmt"
)

// Log 项目基础配置
type Log struct {
	Path      string `json:"path" yaml:"path"`
	AccessLog bool   `json:"access_log" yaml:"access_log"`
}

func (l Log) LogFilePath(filename string) string {
	return fmt.Sprintf("%s/logs/%s", l.Path, filename)
}
