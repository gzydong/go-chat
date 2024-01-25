package socket

type IConn interface {
	// Read 数据读取
	Read() ([]byte, error)
	// Write 数据写入
	Write([]byte) error
	// Close 连接关闭
	Close() error
	// SetCloseHandler 设置连接关闭回调事件
	SetCloseHandler(fn func(code int, text string) error)
	// Network 网络协议类型
	Network() string
}
