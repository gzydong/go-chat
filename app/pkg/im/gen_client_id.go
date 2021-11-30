package im

import "sync/atomic"

// 全局发号器
var GenClientID = &genClientID{}

type genClientID struct {
	number uint64 // 当前号码牌
}

// GetID 获取自增ID
func (gen *genClientID) GetID() int64 {
	return int64(atomic.AddUint64(&gen.number, 1))
}

// GetMaxID 获取当前最大自增ID
func (gen *genClientID) GetMaxID() int64 {
	return int64(atomic.LoadUint64(&gen.number))
}
