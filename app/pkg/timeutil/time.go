package timeutil

import "time"

const (
	DatetimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
)

func DateTime() string {
	return time.Now().Format(DatetimeFormat)
}

func Date() string {
	return time.Now().Format(DateFormat)
}

func Timestamp() int64 {
	return time.Now().Unix()
}

func Time() string {
	return time.Now().Format(TimeFormat)
}
