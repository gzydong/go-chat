package worker

import (
	"fmt"
	"testing"
	"time"
)

func TestConsume_Start(t *testing.T) {
	c := NewConsume(1, 2)

	c.Start()

	for i := 0; i < 100; i++ {
		s := i
		c.AddTask(s, TaskFun(func(index int) {
			fmt.Println("RUN ", index, "-", s)
			time.Sleep(1 * time.Second)
		}))
	}

	c.StopTask()

	c.Wait()
}
