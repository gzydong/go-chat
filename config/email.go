package config

type Email struct {
	Driver      string `yaml:"driver"`       // 邮件驱动
	Host        string `yaml:"host"`         // smtp.163.com
	Port        int    `yaml:"port"`         // 端口号
	UserName    string `yaml:"user_name"`    // 登录账号
	Password    string `yaml:"password"`     // 登录密码
	FromAddress string `yaml:"from_address"` // 邮件地址
	FromName    string `yaml:"from_name"`    // 发送人名称
}
