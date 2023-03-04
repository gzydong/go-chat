package timewheel

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"go-chat/internal/pkg/timeutil"
)

func TestNewSimpleTimeWheel(t *testing.T) {

	var num int32

	obj := NewSimpleTimeWheel(1*time.Second, 100, func(wheel *SimpleTimeWheel, value any) {

		atomic.AddInt32(&num, 1)

		if val, ok := value.(*Conn); ok {
			fmt.Println("预期过期时间", val.lastTime.Format(timeutil.DatetimeFormat), "当前时间", time.Now().Format(timeutil.DatetimeFormat), "num", atomic.LoadInt32(&num))
		}
	})

	go obj.Start()

	for i := 0; i < 10000; i++ {

		go func() {
			cn := &Conn{lastTime: time.Now().Add(time.Duration(13) * time.Second)}
			obj.Add(cn, time.Duration(13)*time.Second)

			time.Sleep(1 * time.Second)
			obj.Remove(cn)
		}()

	}

	time.Sleep(1 * time.Hour)
}
