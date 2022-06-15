package other

import (
	"context"
	"fmt"
	"time"

	"go-chat/internal/pkg/worker"
)

type ExampleHandle struct {
}

func NewExampleHandle() *ExampleHandle {
	return &ExampleHandle{}
}

func (e *ExampleHandle) Handle(ctx context.Context) error {
	c := worker.NewConsume(10, 2)

	for i := 0; i < 100; i++ {
		s := i

		c.AddTask(s, worker.TaskFunc(func(index int) {
			fmt.Println("RUN ", index, "-", s)
			time.Sleep(3 * time.Second)
		}))
	}

	c.StopTask()

	c.Wait()

	fmt.Println("NewExampleHandle")
	return nil
}
