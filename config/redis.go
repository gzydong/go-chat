package config

// Redis Redis配置信息
type Redis struct {
	Host     string `json:"host" yaml:"host"`         // 服务器IP地址
	Port     int    `json:"port" yaml:"port"`         // 服务器端口号
	Auth     string `json:"auth" yaml:"auth"`         // 服务器端口号
	Database int    `json:"database" yaml:"database"` // 数据库
}
