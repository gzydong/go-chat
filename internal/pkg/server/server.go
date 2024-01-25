package server

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"sync"
)

var (
	once sync.Once
	// 服务唯一ID
	serverId string
)

func init() {
	once.Do(func() {
		id, err := gonanoid.Generate("0123456789abcdefghjklmnpqrstuvwxyz", 10)
		if err != nil {
			panic(err)
		}

		serverId = id
	})
}

func ID() string {
	return serverId
}
