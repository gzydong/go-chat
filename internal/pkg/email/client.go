package email

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Client struct {
	config *MailConfig
}

type MailConfig struct {
	Host     string // smtp.163.com
	Port     int    // 端口号
	UserName string // 登录账号
	Password string // 登录密码
	FromName string // 发送人名称
}

func NewEmail(config *MailConfig) *Client {

	fmt.Println(config)
	return &Client{
		config: config,
	}
}

type OptionFunc func(message *gomail.Message)

type Option struct {
	To      []string
	Subject string
	Body    string
}

func (e *Client) SendMail(email *Option, opt ...OptionFunc) error {
	m := gomail.NewMessage()

	fmt.Println(e.config)

	m.SetHeader("From", m.FormatAddress(e.config.UserName, e.config.FromName)) // 这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", email.To...)                                             // 发送给多个用户
	m.SetHeader("Subject", email.Subject)                                      // 设置邮件主题
	m.SetBody("text/html", email.Body)                                         // 设置邮件正文

	for _, o := range opt {
		o(m)
	}

	dialer := gomail.NewDialer(e.config.Host, e.config.Port, e.config.UserName, e.config.Password)

	// 自动开启SSL
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(m)
}
