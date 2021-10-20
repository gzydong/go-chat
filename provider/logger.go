package provider

import (
	"github.com/sirupsen/logrus"
	"go-chat/config"
	"os"
)

func NewLogger(conf *config.Config) *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetOutput(os.Stdout)

	return log
}
