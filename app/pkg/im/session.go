package im

import (
	"sync"
)

// SessionManage session 客户端管理实例
var SessionManage = &session{
	DefaultChannel: &ChannelManage{
		Name:    "default",
		Count:   0,
		Clients: make(map[int64]*Client),
		Lock:    &sync.RWMutex{},
		inChan:  make(chan *ReceiveContent, 3<<20),
		outChan: make(chan *SenderContent, 3<<20),
	},
}

// session 渠道客户端
type session struct {
	DefaultChannel *ChannelManage // 默认分组
	// 可注册其它渠道...
}
