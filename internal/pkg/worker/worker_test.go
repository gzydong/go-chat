package worker

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	w := NewWorker(100, 100)

	var total int32

	for i := 0; i < 10000; i++ {
		num := i
		w.Do(func() {
			fmt.Println("num:", num)
			time.Sleep(1 * time.Second)

			atomic.AddInt32(&total, 1)
		})
	}

	w.Wait()

	fmt.Println("total: ", total)
}
