package timewheel

import (
	"fmt"
	"testing"
	"time"

	"go-chat/internal/pkg/strutil"
)

func TestNewSimpleTimeWheel(t *testing.T) {

	obj := NewSimpleTimeWheel(1*time.Second, 100, func(wheel *SimpleTimeWheel, value any) {
		fmt.Println(value)
	})

	go obj.Start()

	for i := 0; i < 10000; i++ {

		if i%200 == 0 {
			time.Sleep(1 * time.Second)
		}

		obj.Add(strutil.NewMsgId(), 2*time.Second)
	}

	time.Sleep(1 * time.Hour)
}
