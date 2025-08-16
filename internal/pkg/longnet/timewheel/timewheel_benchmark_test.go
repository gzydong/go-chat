package timewheel

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkTimingWheel_AddAndExecuteTask(b *testing.B) {
	const (
		interval = 10 * time.Millisecond
		slots    = 10
		workers  = 10
		delay    = 25 * time.Millisecond
	)

	tw := NewTimingWheel(interval, slots, workers)
	defer tw.Stop()

	var wg sync.WaitGroup
	wg.Add(b.N)

	// 设置回调函数
	tw.SetCallback(func(taskId int64) {
		wg.Done()
	})

	// 开始计时
	b.ResetTimer()

	// 并发添加任务
	for i := 0; i < b.N; i++ {
		err := tw.AddTask(int64(i), delay, nil)
		if err != nil {
			b.Fatalf("Failed to add task %d: %v", i, err)
		}
	}
	// 等待所有任务完成
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		b.Fatal("Timeout waiting for tasks to complete")
	}

	// 停止计时
	b.StopTimer()
}
