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
			RecvChan: make(chan []byte, 1000),
			SendChan: make(chan *Message, 1000),
		},
		//AdminChannel: &ChannelManager{
		//	Name:     "admin",
		//	Count:    0,
		//	Clients:  make(map[string]*Client),
		//	Lock:     &sync.Mutex{},
		//	RecvChan: make(chan []byte, 0),
		//	SendChan: make(chan *Message, 0),
		//},

		// 可注册其它渠道...
	}
}

// StartServer 启动协程处理推送信息
func StartServer() {
	// 开启一个协程消费通道信息
	// (注)添加新的业务渠道续手动添加消费协程
	go Manager.DefaultChannel.ConsumerProcess()

	// 暂时用不到预留
	//go Manager.AdminChannel.ConsumerProcess()
}
