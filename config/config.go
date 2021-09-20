package config

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GlobalConfig *Config

// GlobalConfig 为系统全局配置
type Config struct {
	Redis  Redis `json:"redis" yaml:"redis"`
	MySQL  MySQL `json:"mysql" yaml:"mysql"`
	Jwt    Jwt   `json:"jwt" yaml:"jwt"`
	Server Server
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

	// 生成服务运行ID
	GlobalConfig.Server.RunID = uuid.NewV4().String()

	//fmt.Printf("config %#v\n", GlobalConfig)

	fmt.Println("ServerID:", GetServerRunId())
}

// GetServerRunId 获取当前服务运行ID(服务重启后会改变)
func GetServerRunId() string {
	return GlobalConfig.Server.RunID
}
