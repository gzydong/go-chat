package im

import (
	"sync"
)

// Manager 客户端管理实例
var Manager = &ChannelGroup{
	DefaultChannel: &ChannelManager{
		Name:    "default",
		Count:   0,
		Clients: make(map[int]*Client),
		Lock:    &sync.RWMutex{},
		inChan:  make(chan *RecvMessage, 10240),
		outChan: make(chan *SendMessage, 10240),
	},
}

// ChannelGroup 渠道客户端
type ChannelGroup struct {
	DefaultChannel *ChannelManager // 默认分组
	// 可注册其它渠道...
}
