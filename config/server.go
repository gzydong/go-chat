package config

// Server 项目基础配置
type Server struct {
	AppName  string `json:"app_name" yaml:"app_name"` // 项目名称
	Port     int    `json:"port" yaml:"port"`         // http 启动端口号
	Version  string `json:"version" yaml:"version"`   // 版本号
	ServerId string `json:"server_id"`                // 服务运行ID（程序启动自动生成，每次生成唯一）
	Debug    bool   `json:"debug" yaml:"debug"`
}
