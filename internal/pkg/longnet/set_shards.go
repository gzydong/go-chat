package longnet

import (
	"sync"
)

// 定义分段数量（通常为 16 或 32）
const shardCount = 16

type shardsItem struct {
	mu    sync.RWMutex
	items map[int64][]int64
}

type SetShards struct {
	shards []*shardsItem
}

func NewSetShards() *SetShards {
	m := &SetShards{
		shards: make([]*shardsItem, 0, shardCount),
	}

	for i := 0; i < shardCount; i++ {
		m.shards = append(m.shards, &shardsItem{
			items: make(map[int64][]int64, 32),
		})
	}

	return m
}

func (s *SetShards) getShard(uid int64) *shardsItem {
	return s.shards[fnv32(uid)%uint32(shardCount)]
}

// Get 获取某个用户对应的所有 sid
func (s *SetShards) Get(uid int64) []int64 {
	shard := s.getShard(uid)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	return shard.items[uid]
}

// Add 为用户添加一个 sid
func (s *SetShards) Add(uid int64, sid int64) {
	shard := s.getShard(uid)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	sids, ok := shard.items[uid]
	if !ok {
		sids = []int64{sid}
	} else {
		for _, v := range sids {
			if v == sid {
				return
			}
		}
		sids = append(sids, sid)
	}

	shard.items[uid] = sids
}

// Del 删除用户的一个 sid
func (s *SetShards) Del(uid int64, sid int64) {
	shard := s.getShard(uid)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	sids, ok := shard.items[uid]
	if !ok {
		return
	}

	newSids := make([]int64, 0, len(sids))
	for _, v := range sids {
		if v != sid {
			newSids = append(newSids, v)
		}
	}

	if len(newSids) == 0 {
		delete(shard.items, uid)
	} else {
		shard.items[uid] = newSids
	}
}

func (s *SetShards) GetUserNum() int32 {

	num := 0
	for _, shard := range s.shards {
		num += len(shard.items)
	}

	return int32(num)
}
