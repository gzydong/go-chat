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

func StartServer() {
	// 开启一个协程消费通道信息
	// (注)添加新的业务渠道续手动添加消费协程
	go Manager.DefaultChannel.ConsumerProcess()
}
