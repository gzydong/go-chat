package socket

import (
	"sync"
)

var _ IRoomStorage = (*RoomStorage)(nil)

// IRoomStorage 定义房间存储接口
type IRoomStorage interface {
	// Insert 添加房间成员
	Insert(groupId int32, clientId int64, timestamp int64) error
	// BatchInsert 批量添加房间成员
	BatchInsert(groupId int32, clientIds []int64, timestamp int64) error
	// Delete 删除房间成员
	Delete(groupId int32, clientId int64, timestamp int64) error
	// BatchDelete 批量删除房间成员
	BatchDelete(groupId int32, clientIds []int64, timestamp int64) error
	// IsRoomMember 判断用户是否在房间内
	IsRoomMember(groupId int32, clientId int64) bool
	// GetClientIDAll 获取房间所有成员
	GetClientIDAll(groupId int32) []int64
	// DeleteRoom 删除房间
	DeleteRoom(groupId int32) error
	// GetRoomNum 获取房间数量
	GetRoomNum() int32
}

// RoomEntity 表示一个房间及其成员
type RoomEntity struct {
	mutex sync.RWMutex
	items map[int64]int64
}

// RoomStorage 实现了房间存储
type RoomStorage struct {
	mutex sync.RWMutex
	rooms map[int32]*RoomEntity
}

// NewRoomStorage 创建一个新的RoomStorage实例
func NewRoomStorage() *RoomStorage {
	return &RoomStorage{
		rooms: make(map[int32]*RoomEntity, 1000),
	}
}

func (r *RoomStorage) GetRoomNum() int32 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return int32(len(r.rooms))
}

func (r *RoomStorage) Insert(groupId int32, clientId int64, timestamp int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	entity, ok := r.rooms[groupId]
	if !ok {
		entity = &RoomEntity{
			mutex: sync.RWMutex{},
			items: make(map[int64]int64),
		}
		r.rooms[groupId] = entity
	}

	entity.mutex.Lock()
	entity.items[clientId] = timestamp
	entity.mutex.Unlock()

	return nil
}

func (r *RoomStorage) BatchInsert(groupId int32, clientIds []int64, timestamp int64) error {
	r.mutex.Lock()

	entity, ok := r.rooms[groupId]
	if !ok {
		entity = &RoomEntity{
			mutex: sync.RWMutex{},
			items: make(map[int64]int64, len(clientIds)),
		}

		for _, id := range clientIds {
			entity.items[id] = timestamp
		}

		r.rooms[groupId] = entity
		r.mutex.Unlock()
	} else {
		r.mutex.Unlock()

		entity.mutex.Lock()
		for _, id := range clientIds {
			entity.items[id] = timestamp
		}

		entity.mutex.Unlock()
	}

	return nil
}

func (r *RoomStorage) Delete(groupId int32, clientId int64, timestamp int64) error {
	r.mutex.RLock()
	entity, ok := r.rooms[groupId]
	r.mutex.RUnlock()
	if !ok {
		return nil
	}

	entity.mutex.Lock()
	defer entity.mutex.Unlock()

	if value, ok := entity.items[clientId]; ok && value < timestamp {
		delete(entity.items, clientId)
	}

	return nil
}

func (r *RoomStorage) BatchDelete(groupId int32, clientIds []int64, timestamp int64) error {
	r.mutex.RLock()
	entity, ok := r.rooms[groupId]
	r.mutex.RUnlock()

	if !ok {
		return nil
	}

	entity.mutex.Lock()
	defer entity.mutex.Unlock()

	for _, id := range clientIds {
		if value, ok := entity.items[id]; ok && value < timestamp {
			delete(entity.items, id)
		}
	}

	return nil
}

func (r *RoomStorage) IsRoomMember(groupId int32, clientId int64) bool {
	r.mutex.RLock()
	entity, ok := r.rooms[groupId]
	r.mutex.RUnlock()
	if !ok {
		return false
	}

	entity.mutex.RLock()
	defer entity.mutex.RUnlock()

	_, ok = entity.items[clientId]
	return ok
}

func (r *RoomStorage) GetClientIDAll(groupId int32) []int64 {
	r.mutex.RLock()
	entity, ok := r.rooms[groupId]
	r.mutex.RUnlock()

	if !ok {
		return make([]int64, 0)
	}

	entity.mutex.RLock()
	defer entity.mutex.RUnlock()

	uids := make([]int64, 0, len(entity.items))
	for uid := range entity.items {
		uids = append(uids, uid)
	}

	return uids
}

func (r *RoomStorage) DeleteRoom(groupId int32) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	entity, ok := r.rooms[groupId]
	if !ok {
		return nil
	}

	delete(r.rooms, groupId)

	entity.mutex.Lock()
	entity.items = make(map[int64]int64) // 清空房间成员
	entity.mutex.Unlock()

	return nil
}
