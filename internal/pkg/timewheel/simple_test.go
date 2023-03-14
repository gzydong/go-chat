package timewheel

import (
	"fmt"
	"testing"
	"time"
)

func TestNewSimpleTimeWheel(t *testing.T) {

	obj := NewSimpleTimeWheel(1*time.Second, 100, func(wheel *SimpleTimeWheel, key string, value any) {
		fmt.Println(value)
	})

	go obj.Start()

	time.Sleep(1 * time.Hour)
}
