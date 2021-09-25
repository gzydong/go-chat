package config

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GlobalConfig *Config

// GlobalConfig 为系统全局配置
type Config struct {
	Redis  Redis  `json:"redis" yaml:"redis"`
	MySQL  MySQL  `json:"mysql" yaml:"mysql"`
	Jwt    Jwt    `json:"jwt" yaml:"jwt"`
	Cors   Cors   `json:"cors" yaml:"cors"`
	Server Server `json:"server" yaml:"server"`
}

// NewConfig
func NewConfig() {
	GlobalConfig = &Config{}

	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	if yaml.Unmarshal(content, &GlobalConfig) != nil {
		panic(fmt.Sprintf("解析 config.yaml 读取错误: %v", err))
	}

	// 生成服务运行ID
	GlobalConfig.Server.ServerID = uuid.NewV4().String()

	text, _ := jsoniter.MarshalToString(GlobalConfig)
	fmt.Printf("项目配置信息: %s\n\n", text)
}

// GetServerID 获取当前服务运行ID(服务重启后会改变)
func GetServerID() string {
	return GlobalConfig.Server.ServerID
}
