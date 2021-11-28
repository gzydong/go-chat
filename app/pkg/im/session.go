package im

import (
	"sync"
)

// Sessions 客户端管理实例
var Sessions = &session{
	Default: &ChannelManage{
		Name:    "default",
		Count:   0,
		Clients: make(map[int64]*Client),
		Lock:    &sync.RWMutex{},
		inChan:  make(chan *ReceiveContent, 5<<20),
		outChan: make(chan *SenderContent, 5<<20),
	},
}

// session 渠道客户端
type session struct {
	Default *ChannelManage // 默认分组
	// 可注册其它渠道...
}
