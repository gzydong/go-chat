package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 配置信息
type Config struct {
	App        *App        `json:"app" yaml:"app"`
	Redis      *Redis      `json:"redis" yaml:"redis"`
	MySQL      *MySQL      `json:"mysql" yaml:"mysql"`
	Jwt        *Jwt        `json:"jwt" yaml:"jwt"`
	Cors       *Cors       `json:"cors" yaml:"cors"`
	Log        *Log        `json:"log" yaml:"log"`
	Filesystem *Filesystem `json:"filesystem" yaml:"filesystem"`
	Email      *Email      `json:"email" yaml:"email"`
	Server     *Server     `json:"server" yaml:"server"`
	Nsq        *Nsq        `json:"nsq" yaml:"nsq"` // 目前没用到
	OAuth      *OAuth      `json:"oauth" yaml:"oauth"`
}

type Server struct {
	HttpAddr      string `json:"http_addr" yaml:"http_addr"`
	WebsocketAddr string `json:"websocket_addr" yaml:"websocket_addr"`
	TcpAddr       string `json:"tcp_addr" yaml:"tcp_addr"`
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

	return &conf
}

// Debug 调试模式
func (c *Config) Debug() bool {
	return c.App.Debug
}
