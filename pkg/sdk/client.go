package sdk

import (
	"net"
)

type ITcpClient interface {
	Connect() error
	// Send 推送消息
	Send() error
	// Message 消息回调
	Message(data []byte) error
	// Close 关闭连接
	Close() error
}

// nolint
type TcpClient struct {
	// 服务端地址
	address string
	// tcp 连接
	conn net.Conn
}

func NewTcpClient() *TcpClient {
	return &TcpClient{}
}
