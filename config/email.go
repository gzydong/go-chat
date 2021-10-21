package config

// Email 邮件配置信息
type Email struct {
	Host     string `yaml:"host"`     // smtp.163.com
	Port     int    `yaml:"port"`     // 端口号
	UserName string `yaml:"username"` // 登录账号
	Password string `yaml:"password"` // 登录密码
	FromName string `yaml:"fromname"` // 发送人名称
}
