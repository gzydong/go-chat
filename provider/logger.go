package provider

import (
	"github.com/sirupsen/logrus"
	"go-chat/config"
)

func NewLogger(conf *config.Config) *logrus.Logger {
	log := logrus.New()

	return log
}
