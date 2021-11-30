package im

// Sessions 客户端管理实例
var Sessions = &session{
	// Default: &Channel{
	// 	Name:    "default",
	// 	count:   0,
	// 	inChan:  make(chan *ReceiveContent, 5<<20),
	// 	outChan: make(chan *SenderContent, 5<<20),
	// 	maps:    maps(10), // 拆分 map 数，可合理分配
	// },

	Default: NewChannel("default", maps(10), make(chan *ReceiveContent, 5<<20), make(chan *SenderContent, 5<<20)),
}

// session 渠道客户端
type session struct {
	Default *Channel // 默认分组
	// 可注册其它渠道...
}
