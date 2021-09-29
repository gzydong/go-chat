package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

// Config 配置信息
type Config struct {
	Redis      RedisConfig `json:"redis" yaml:"redis"`
	MySQL      MySQL       `json:"mysql" yaml:"mysql"`
	Jwt        Jwt         `json:"jwt" yaml:"jwt"`
	Cors       Cors        `json:"cors" yaml:"cors"`
	Server     Server      `json:"server" yaml:"server"`
	Filesystem Filesystem  `json:"filesystem" yaml:"filesystem"`
}

func Init(filename string) *Config {
	conf := &Config{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if yaml.Unmarshal(content, conf) != nil {
		panic(fmt.Sprintf("解析config.yaml读取错误: %v", err))
	}

	// 生成服务运行ID
	conf.Server.ServerId = strings.Replace(uuid.NewV4().String(), "-", "", 4)

	return conf
}
