package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"time"

	"github.com/natefinch/lumberjack"
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

func CreateFileWriter(filePath string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filePath, // 日志文件的位置
		MaxSize:    100,      // 文件最大尺寸（以MB为单位）
		MaxBackups: 3,        // 保留的最大旧文件数量
		MaxAge:     7,        // 保留旧文件的最大天数
		Compress:   true,     // 是否压缩/归档旧文件
		LocalTime:  true,     // 使用本地时间创建时间戳
	}
}

func Init(filePath string, level slog.Level, topic string) {
	handler := slog.NewJSONHandler(CreateFileWriter(filePath), &slog.HandlerOptions{
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
	logf(context.Background(), slog.LevelDebug, format, args...)
}

func Debug(msg string) {
	logf(context.Background(), slog.LevelDebug, msg)
}

func Infof(format string, args ...any) {
	logf(context.Background(), slog.LevelInfo, format, args...)
}

func Info(msg string) {
	logf(context.Background(), slog.LevelInfo, msg)
}

func Warnf(format string, args ...any) {
	logf(context.Background(), slog.LevelWarn, format, args...)
}

func Warn(msg string) {
	logf(context.Background(), slog.LevelWarn, msg)
}

func Errorf(format string, args ...any) {
	logf(context.Background(), slog.LevelError, format, args...)
}

func Error(msg string) {
	logf(context.Background(), slog.LevelError, msg)
}

func logf(ctx context.Context, level slog.Level, format string, args ...any) {
	if !out.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}

	r := slog.NewRecord(time.Now(), level, format, pcs[0])

	_ = out.Handler().Handle(ctx, r)
}

func WithFields(level slog.Level, msg string, fields any) {
	if !out.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add("extra", fields)

	_ = out.Handler().Handle(context.Background(), r)
}

func ErrorWithFields(msg string, err error, fields any) {
	ctx := context.Background()

	if !out.Enabled(ctx, slog.LevelError) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])

	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	r.Add("error", err.Error())
	r.Add("extra", fields)

	_ = out.Handler().Handle(ctx, r)
}
