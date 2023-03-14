package adapter

import (
	"fmt"
	"io"
	"net"

	"go-chat/internal/pkg/ichat/socket/adapter/encoding"
)

// TcpAdapter TCP 适配器
type TcpAdapter struct {
	conn      net.Conn
	hookClose func(code int, text string) error
}

func NewTcpAdapter(conn net.Conn) (*TcpAdapter, error) {
	return &TcpAdapter{conn: conn}, nil
}

func (t *TcpAdapter) Network() string {
	return TcpType
}

func (t *TcpAdapter) Read() ([]byte, error) {

	msg, err := encoding.NewDecode(t.conn)

	if err == io.EOF {
		if t.hookClose != nil {
			if err := t.hookClose(1000, "客户端已关闭"); err != nil {
				return nil, err
			}
		}

		return nil, fmt.Errorf("连接已断开")
	}

	if err != nil {
		return nil, fmt.Errorf("decode msg failed, err: %s", err.Error())
	}

	return msg, nil
}

func (t *TcpAdapter) Write(bytes []byte) error {

	binaryData, err := encoding.NewEncode(bytes)
	if err != nil {
		return err
	}

	_, err = t.conn.Write(binaryData)

	return err
}

func (t *TcpAdapter) Close() error {
	return t.conn.Close()
}

func (t *TcpAdapter) SetCloseHandler(fn func(code int, text string) error) {
	t.hookClose = fn
}
