package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/strutil"
	"gopkg.in/yaml.v2"
)

// Config 配置信息
type Config struct {
	sid        string      // 服务运行ID
	App        *App        `json:"app" yaml:"app"`
	Redis      *Redis      `json:"redis" yaml:"redis"`
	MySQL      *MySQL      `json:"mysql" yaml:"mysql"`
	Jwt        *Jwt        `json:"jwt" yaml:"jwt"`
	Cors       *Cors       `json:"cors" yaml:"cors"`
	Log        *Log        `json:"log" yaml:"log"`
	Filesystem *Filesystem `json:"filesystem" yaml:"filesystem"`
	Email      *Email      `json:"email" yaml:"email"`
}

func ReadConfig(filename string) *Config {
	conf := &Config{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if yaml.Unmarshal(content, conf) != nil {
		panic(fmt.Sprintf("解析 config.yaml 读取错误: %v", err))
	}

	// 生成服务运行ID
	conf.sid = encrypt.Md5(fmt.Sprintf("%d%s", time.Now().UnixNano(), strutil.Random(6)))

	return conf
}

// ServerId 服务运行ID
func (c *Config) ServerId() string {
	return c.sid
}

// Debug 调试模式
func (c *Config) Debug() bool {
	return c.App.Debug
}

func (c *Config) SetPort(port int) {
	c.App.Port = port
}

func (c *Config) GetLogPath() string {
	return c.Log.Path
}
