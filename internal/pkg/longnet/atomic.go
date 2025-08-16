package longnet

import (
	"sync/atomic"
)

type AtomicInt64 struct {
	value int64
}

// Incr 自增+1
func (a *AtomicInt64) Incr() int64 {
	return atomic.AddInt64(&a.value, 1)
}

// Decr 自减-1
func (a *AtomicInt64) Decr() int64 {
	return atomic.AddInt64(&a.value, -1)
}

func (a *AtomicInt64) Load() int64 {
	return atomic.LoadInt64(&a.value)
}

type AtomicInt32 struct {
	value int32
}

// Incr 自增+1
func (a *AtomicInt32) Incr() int32 {
	return atomic.AddInt32(&a.value, 1)
}

// Decr 自减-1
func (a *AtomicInt32) Decr() int32 {
	return atomic.AddInt32(&a.value, -1)
}

func (a *AtomicInt32) Load() int32 {
	return atomic.LoadInt32(&a.value)
}
