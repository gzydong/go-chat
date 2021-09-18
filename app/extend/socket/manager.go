package socket

import (
	"reflect"
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
			ChannelName: "default",
			Count:       0,
			Clients:     make(map[string]*WsClient),
			Lock:        &sync.Mutex{},
		},

		// 可注册其它渠道...
	}
}

// StartServer 启动服务
func StartServer() {
	el := reflect.ValueOf(Manager).Elem()
	for i := 0; i < el.NumField(); i++ {

	}
}
