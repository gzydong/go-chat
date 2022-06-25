package provider

import (
	"go-chat/config"
	"go-chat/internal/pkg/email"
)

func NewEmailClient(conf *config.Config) *email.Client {
	return email.NewEmail(&email.Config{
		Host:     conf.Email.Host,
		Port:     conf.Email.Port,
		UserName: conf.Email.UserName,
		Password: conf.Email.Password,
		FromName: conf.Email.FromName,
	})
}
