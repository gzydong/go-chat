package worker

import (
	"fmt"
	"testing"
	"time"
)

func TestConsume_Start(t *testing.T) {
	c := NewConsume(50, 10)

	for i := 0; i < 1000; i++ {
		s := i

		c.AddTask(s, TaskFunc(func(index int) {
			fmt.Println("RUN ", index, "-", s)
			time.Sleep(1 * time.Second)
		}))
	}

	c.StopTask()
	c.Wait()
}
