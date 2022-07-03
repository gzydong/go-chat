package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Client struct {
	config *Config
}

type Config struct {
	Host     string // 例如 smtp.163.com
	Port     int    // 端口号
	UserName string // 登录账号
	Password string // 登录密码
	FromName string // 发送人名称
}

func NewEmail(config *Config) *Client {
	return &Client{
		config: config,
	}
}

type Option struct {
	To      []string // 收件人
	Subject string   // 邮件主题
	Body    string   // 邮件正文
}

type OptionFunc func(msg *gomail.Message)

func (c *Client) do(msg *gomail.Message) error {
	dialer := gomail.NewDialer(c.config.Host, c.config.Port, c.config.UserName, c.config.Password)

	// 自动开启SSL
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(msg)
}

func (c *Client) SendMail(email *Option, opt ...OptionFunc) error {
	m := gomail.NewMessage()

	// 这种方式可以添加别名，即“XX官方”
	m.SetHeader("From", m.FormatAddress(c.config.UserName, c.config.FromName))

	if len(email.To) > 0 {
		m.SetHeader("To", email.To...)
	}

	if len(email.Subject) > 0 {
		m.SetHeader("Subject", email.Subject)
	}

	if len(email.Body) > 0 {
		m.SetBody("text/html", email.Body)
	}

	// m.SetHeader("Cc", m.FormatAddress("xxxx@foxmail.com", "收件人")) //抄送
	// m.SetHeader("Bcc", m.FormatAddress("xxxx@gmail.com", "收件人"))  //暗送

	for _, o := range opt {
		o(m)
	}

	return c.do(m)
}
