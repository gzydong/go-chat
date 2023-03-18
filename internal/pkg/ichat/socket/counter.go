package socket

import "sync/atomic"

// Counter 全局发号器
var Counter = &counter{}

type counter struct {
	number uint64
}

// GenID 获取自增ID
func (c *counter) GenID() int64 {
	return int64(atomic.AddUint64(&c.number, 1))
}

// GetMaxID 获取当前最大自增ID
func (c *counter) GetMaxID() int64 {
	return int64(atomic.LoadUint64(&c.number))
}
