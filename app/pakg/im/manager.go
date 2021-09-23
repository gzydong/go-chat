package im

import (
	"sync"
)

// 客户端管理实例
var Manager *ChannelGroup

// 渠道客户端
type ChannelGroup struct {
	DefaultChannel *ChannelManager // 默认分组
	AdminChannel   *ChannelManager // 后台分组
	// 可注册其它渠道...
}

// init 初始化注册分组
func init() {
	Manager = &ChannelGroup{
		DefaultChannel: &ChannelManager{
			Name:     "default",
			Count:    0,
			Clients:  make(map[string]*Client),
			Lock:     &sync.Mutex{},
			RecvChan: make(chan *RecvMessage, 10000),
			SendChan: make(chan *SendMessage, 10000),
		},
		AdminChannel: &ChannelManager{
			Name:     "admin",
			Count:    0,
			Clients:  make(map[string]*Client),
			Lock:     &sync.Mutex{},
			RecvChan: make(chan *RecvMessage, 0),
			SendChan: make(chan *SendMessage, 0),
		},

		// 可注册其它渠道...
	}
}
