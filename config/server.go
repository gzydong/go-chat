package config

type Server struct {
	AppName  string `json:"app_name" yaml:"app_name"` // 项目名称
	Version  string `json:"version" yaml:"version"`   // 版本号
	ServerID string `json:"server_id"`                // 服务运行ID（程序启动自动生成，每次生成唯一）
}
