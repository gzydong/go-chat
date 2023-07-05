package config

import "fmt"

// MySQL 数据库配置信息
type MySQL struct {
	Host      string `json:"host" yaml:"host"`         // 服务器IP地址
	Port      int    `json:"port" yaml:"port"`         // 服务器端口号
	UserName  string `json:"username" yaml:"username"` // 数据库用户名
	Password  string `json:"password" yaml:"password"` // 数据库用户密码
	Database  string `json:"database" yaml:"database"` // 数据库名
	Prefix    string `json:"prefix" yaml:"prefix"`     // 数据表前缀
	Charset   string `json:"charset" yaml:"charset"`
	Collation string `json:"collation" yaml:"collation"`
}

func (m *MySQL) Dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.UserName,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
		m.Charset,
	)
}
