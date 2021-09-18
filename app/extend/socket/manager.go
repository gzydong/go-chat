package socket

import (
	"sync"
)

// 客户端管理实例
var Manager *ChannelGroup

// 渠道客户端
type ChannelGroup struct {
	DefaultChannel *ChannelManager // 默认分组

	// 可注册其它渠道...
}

// 初始化注册分组
func init() {
	Manager = &ChannelGroup{
		DefaultChannel: &ChannelManager{
			Name:        "default",
			Count:       0,
			Clients:     make(map[string]*Client),
			Lock:        &sync.Mutex{},
			ChanMessage: make(chan *Message, 1000),
		},

		// 可注册其它渠道...
	}
}
