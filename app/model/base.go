package model

import (
	"go-chat/app/pkg/timeutil"
	"time"
)

type DateTime time.Time

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeutil.DatetimeFormat+`"`, string(data), time.Local)
	*t = DateTime(now)
	return
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeutil.DatetimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeutil.DatetimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t DateTime) String() string {
	return time.Time(t).Format(timeutil.DatetimeFormat)
}
