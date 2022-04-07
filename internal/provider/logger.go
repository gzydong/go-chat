package provider

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	_ = os.MkdirAll("./runtime/logs", os.ModePerm)

	src, err := os.OpenFile("./runtime/logs/logger.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(src)
}
