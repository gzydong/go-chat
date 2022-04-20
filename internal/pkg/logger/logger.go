package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var out *logrus.Logger

func init() {
	out = logrus.New()

	out.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func SetOutput(dir string, name string) {
	_ = os.MkdirAll(dir, os.ModePerm)

	if !strings.HasSuffix(name, ".log") {
		name = fmt.Sprintf("%s.log", name)
	}

	src, err := os.OpenFile(fmt.Sprintf("%s/logs/%s", dir, name), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	out.SetOutput(src)
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return out.WithFields(logrus.Fields(fields))
}

func Info(args ...interface{}) {
	out.Info(args...)
}

func Infoln(args ...interface{}) {
	out.Infoln(args...)
}

func Infof(format string, args ...interface{}) {
	out.Infof(format, args...)
}

func Warn(args ...interface{}) {
	out.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	out.Warnf(format, args...)
}

func Debug(args ...interface{}) {
	out.Debug(args...)
}

func Error(args ...interface{}) {
	out.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	out.Errorf(format, args...)
}

func Panic(args ...interface{}) {
	out.Panic(args...)
}
