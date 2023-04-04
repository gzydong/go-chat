package config

import (
	"fmt"
	"os"
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
	Server     *Server     `json:"server" yaml:"server"`
	Nsq        *Nsq        `json:"nsq" yaml:"nsq"`
}

type Server struct {
	Http      int `json:"http" yaml:"http"`
	Websocket int `json:"websocket" yaml:"websocket"`
	Tcp       int `json:"tcp" yaml:"tcp"`
}

func New(filename string) *Config {

	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var conf Config
	if yaml.Unmarshal(content, &conf) != nil {
		panic(fmt.Sprintf("解析 config.yaml 读取错误: %v", err))
	}

	// 生成服务运行ID
	conf.sid = encrypt.Md5(fmt.Sprintf("%d%s", time.Now().UnixNano(), strutil.Random(6)))

	return &conf
}

// ServerId 服务运行ID
func (c *Config) ServerId() string {
	return c.sid
}

// Debug 调试模式
func (c *Config) Debug() bool {
	return c.App.Debug
}

func (c *Config) LogPath() string {
	return c.Log.Path
}
