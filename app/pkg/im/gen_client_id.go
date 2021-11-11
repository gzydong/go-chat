package im

import "sync/atomic"

var GenClientID = &genClientID{}

type genClientID struct {
	number uint64 // 当前号码牌
}

func (gen *genClientID) GetID() int64 {
	return int64(atomic.AddUint64(&gen.number, 1))
}

func (gen *genClientID) GetMaxID() int64 {
	return int64(atomic.LoadUint64(&gen.number))
}
