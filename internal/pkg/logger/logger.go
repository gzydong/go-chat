package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	LevelDebug slog.Level = -4
	LevelInfo  slog.Level = 0
	LevelWarn  slog.Level = 4
	LevelError slog.Level = 8
)

var out = slog.Default()

func Std() *slog.Logger {
	return slog.Default()
}

func InitLogger(filePath string, level slog.Level, topic string) {
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		panic(err)
	}

	src, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(src, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}

			return a
		},
	})

	out = slog.New(handler).With("topic", topic)
}

func Debugf(format string, args ...any) {
	logf(slog.LevelDebug, format, args...)
}

func Debug(msg string, args ...any) {
	log(slog.LevelDebug, msg, args...)
}

func Infof(format string, args ...any) {
	logf(slog.LevelInfo, format, args...)
}

func Info(msg string, args ...any) {
	log(slog.LevelInfo, msg, args...)
}

func Warnf(format string, args ...any) {
	logf(slog.LevelWarn, format, args...)
}

func Warn(msg string, args ...any) {
	log(slog.LevelWarn, msg, args...)
}

func Errorf(format string, args ...any) {
	logf(slog.LevelError, format, args...)
}

func Error(msg string, args ...any) {
	log(slog.LevelError, msg, args...)
}

func logf(level slog.Level, format string, args ...any) {
	if !out.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])
	_ = out.Handler().Handle(context.Background(), r)
}

func log(level slog.Level, msg string, args ...any) {
	if !out.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	_ = out.Handler().Handle(context.Background(), r)
}
