package im

// Sessions 客户端管理实例
var Sessions = &session{
	Default: NewChannel("default", NewNode(10), make(chan *ReceiveContent, 5<<20), make(chan *SenderContent, 5<<20)),
}

// session 渠道客户端
type session struct {
	Default *Channel // 默认分组
	// 可注册其它渠道...
}
