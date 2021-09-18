package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GlobalConfig *Config

// GlobalConfig 为系统全局配置
type Config struct {
	Redis Redis `json:"redis" yaml:"redis"`
	MySQL MySQL `json:"mysql" yaml:"mysql"`
	Jwt   Jwt   `json:"jwt" yaml:"jwt"`
}

func init() {
	GlobalConfig = &Config{}
	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	if yaml.Unmarshal(content, &GlobalConfig) != nil {
		panic(fmt.Sprintf("解析config.yaml读取错误: %v", err))
	}

	//fmt.Printf("config %#v\n", GlobalConfig)
}
