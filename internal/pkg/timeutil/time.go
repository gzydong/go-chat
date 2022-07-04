package timeutil

import "time"

const (
	DatetimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	DateDayFormat  = "20060102"
	TimeFormat     = "15:04:05"
)

func DateTime() string {
	return time.Now().Format(DatetimeFormat)
}

func Date() string {
	return time.Now().Format(DateFormat)
}

func DateNumber() string {
	return time.Now().Format(DateDayFormat)
}

func Time() string {
	return time.Now().Format(TimeFormat)
}

func Location() *time.Location {
	lo, _ := time.LoadLocation("Asia/Shanghai")
	return lo
}

func DayStartDateTime() string {
	return Date() + " 00:00:00"
}

func DayEndDateTime() string {
	return Date() + " 23:59:59"
}

func DayStartTime() int64 {
	t, _ := time.ParseInLocation(DatetimeFormat, DayStartDateTime(), Location())

	return t.Unix()
}

func DayEndTime() int64 {
	t, _ := time.ParseInLocation(DatetimeFormat, DayEndDateTime(), Location())
	return t.Unix()
}

// ParseDateTime 解析 Datetime 格式字符串为 Time
func ParseDateTime(datetime string) time.Time {
	t, _ := time.ParseInLocation(DatetimeFormat, datetime, Location())
	return t
}

func FormatDatetime(t time.Time) string {
	return t.Format(DatetimeFormat)
}

func IsDateTimeFormat(datetime string) bool {
	_, err := time.Parse(DatetimeFormat, datetime)

	return err != nil
}

func IsDateFormat(date string) bool {
	_, err := time.Parse(DatetimeFormat, date)

	return err != nil
}
