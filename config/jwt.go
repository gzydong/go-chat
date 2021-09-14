package config

// JWT 相关配置信息
type JWT struct {
	Secret      string `yaml:"jwt-secret"`   // Jwt 秘钥
	ExpiresTime int    `yaml:"expires-time"` // 过期时间(单位秒)
	BufferTime  int    `yaml:"buffer-time"`  // 缓冲时间(单位秒)
}
