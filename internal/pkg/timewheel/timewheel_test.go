package timewheel

import (
	"fmt"
	"testing"
	"time"

	"go-chat/internal/pkg/timeutil"
)

type Conn struct {
	lastTime time.Time
}

// nolint
func TestNewTimeWheel(t *testing.T) {

	var num int

	obj := NewTimeWheel(func(wheel *TimeWheel, a any) {
		num++
		if val, ok := a.(*Conn); ok {
			fmt.Println("预期过期时间", val.lastTime.Format(timeutil.DatetimeFormat), "当前时间", time.Now().Format(timeutil.DatetimeFormat), "num", num)
		}
	})

	go obj.Start()

	tt := time.Now()
	for i := 0; i < 1000; i++ {
		// _ = obj.Add(&Conn{lastTime: time.Now().Add(time.Duration(i) * time.Second)}, time.Duration(i)*time.Second)
		go func(i int) {
			cn := &Conn{lastTime: time.Now().Add(time.Duration(10) * time.Second)}
			obj.Add(cn, time.Duration(10)*time.Second)
			// time.Sleep(3 * time.Second)
			// obj.Remove(cn)
		}(i)

	}

	fmt.Println(time.Since(tt))

	// conn := &Conn{lastTime: time.Now().Add(time.Duration(30) * time.Second)}
	//
	// obj.Add(conn, time.Duration(30)*time.Second)
	// obj.Add(&Conn{lastTime: time.Now().Add(10 * time.Second)}, 10*time.Second)
	// obj.Add(&Conn{lastTime: time.Now().Add(65 * time.Second)}, 65*time.Second)
	// time.Sleep(10 * time.Second)
	// obj.Remove(conn)
	// time.Sleep(4 * time.Second)

	time.Sleep(1 * time.Hour)
}
