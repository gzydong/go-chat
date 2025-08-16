package timewheel

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {

	tw := NewTimingWheel(time.Second, 60, 30)
	tw.SetCallback(func(taskId int64) {

	})

	_ = tw.AddTask(1, time.Second*3, nil)

	<-time.After(5 * time.Second)

	fmt.Println(tw.globalTasks.Get(3))

	tw.Stop()
}
