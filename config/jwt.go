package config

// Jwt 相关配置信息
type Jwt struct {
	Secret      string `yaml:"secret"`       // Jwt 秘钥
	ExpiresTime int64  `yaml:"expires_time"` // 过期时间(单位秒)
}
