package handler

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"go-chat/internal/pkg/im/adapter"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
)

type Handler struct {
	Chat    *ChatChannel
	Example *ExampleChannel
}

type AuthConn struct {
	Uid  int
	conn *adapter.TcpAdapter
}

type Message struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

func (h *Handler) ConnDispatcher(conn net.Conn) {
	ch := make(chan *AuthConn)

	fmt.Println("网络地址", conn.RemoteAddr().(*net.TCPAddr).IP)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover ===>>>", err)
		}
	}()

	go h.auth(conn, ch)

	fmt.Println(conn.RemoteAddr(), "开始认证==>>>", time.Now().Unix())
	select {
	// 认证超时
	case <-time.After(2 * time.Second):
		fmt.Println(conn.RemoteAddr(), "认证超时==>>>", time.Now().Unix())
		_ = conn.Close()
		return

	// 认证成功
	case connInfo := <-ch:
		fmt.Println(conn.RemoteAddr(), "认证成功==>>>", time.Now().Unix())
		fmt.Println(connInfo)

		h.Chat.TcpConn(context.Background(), connInfo.conn)
	}
}

func (*Handler) auth(connect net.Conn, data chan *AuthConn) {
	conn, err := adapter.NewTcpAdapter(connect)
	if err != nil {
		logger.Errorf("tcp connect error: %s", err.Error())
	}

	fmt.Println(connect.RemoteAddr(), "等待认证==>>>", time.Now().Unix())
	read, err := conn.Read()
	if err != nil {
		fmt.Println(connect.RemoteAddr(), "认证异常==>>>", time.Now().Unix(), err.Error())
		return
	}

	msg := &Message{}
	if err := jsonutil.Decode(string(read), msg); err != nil {
		fmt.Println("数据解析失败")
		return
	}

	if msg.Event != "authorize" {
		return
	}

	uid, _ := strconv.Atoi(msg.Content)

	data <- &AuthConn{Uid: uid, conn: conn}
}
