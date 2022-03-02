package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"

	"go-chat/config"
)

type Options struct {
	To      []string
	Subject string
	Body    string
}

func SendMail(config *config.Email, email *Options) error {
	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(config.UserName, config.FromName)) // 这种方式可以添加别名，即“XX官方”

	// 说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	// m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	// m.SetHeader("From", mailConn["user"])

	m.SetHeader("To", email.To...)        // 发送给多个用户
	m.SetHeader("Subject", email.Subject) // 设置邮件主题
	m.SetBody("text/html", email.Body)    // 设置邮件正文

	// m.SetHeader("To", m.FormatAddress("xxxx@qq.com", "收件人"))      // 收件人
	// m.SetHeader("Cc", m.FormatAddress("xxxx@foxmail.com", "收件人")) //抄送
	// m.SetHeader("Bcc", m.FormatAddress("xxxx@gmail.com", "收件人"))  //暗送

	dialer := gomail.NewDialer(config.Host, config.Port, config.UserName, config.Password)

	// 自动开启SSL
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(m)
}
