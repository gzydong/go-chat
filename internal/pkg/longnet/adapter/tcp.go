package adapter

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"go-chat/internal/pkg/longnet/adapter/encoding"
)

//+-------------------+-------------------+-------------------+-------------------+-----------------------------------+
//|      TotalLength  |      Command      |       Flags       |       MsgID       |             Payload              |
//|     (4 bytes)     |     (4 bytes)     |     (4 bytes)     |     (8 bytes)     |        (TotalLength - 20 bytes)   |
//+-------------------+-------------------+-------------------+-------------------+-----------------------------------+

//+-------------------+-------------------------------------------------------+
//|    TotalLength    |                  内层完整消息结构                      |
//|     (4 bytes)     |      (Command + Flags + MsgID + Length + Payload)      |
//+-------------------+-------------------------------------------------------+

// TcpAdapter TCP 适配器
type TcpAdapter struct {
	conn      net.Conn
	reader    *bufio.Reader // Buffer reader for connection.
	hookClose func(code int, text string) error
}

func NewTcpAdapter(conn net.Conn) (*TcpAdapter, error) {
	return &TcpAdapter{conn: conn, reader: bufio.NewReader(conn)}, nil
}

func (t *TcpAdapter) Network() string {
	return NetworkTcp
}

func (t *TcpAdapter) Read() ([]byte, error) {
	msg, err := encoding.NewDecode(t.reader)
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

// SetReadDeadline 设置读取超时时间
func (t *TcpAdapter) SetReadDeadline(deadline time.Time) error {
	return t.conn.SetReadDeadline(deadline)
}

// SetWriteDeadline 设置写入超时时间
func (t *TcpAdapter) SetWriteDeadline(deadline time.Time) error {
	return t.conn.SetWriteDeadline(deadline)
}
