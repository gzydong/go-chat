package im

import (
	"sync"
)

// GroupManage channelGroup 客户端管理实例
var GroupManage = &channelGroup{
	DefaultChannel: &ChannelManager{
		Name:    "default",
		Count:   0,
		Clients: make(map[int]*Client),
		Lock:    &sync.RWMutex{},
		inChan:  make(chan *ReceiveContent, 10240),
		outChan: make(chan *SenderContent, 10240),
	},
}

// channelGroup 渠道客户端
type channelGroup struct {
	DefaultChannel *ChannelManager // 默认分组
	// 可注册其它渠道...
}
