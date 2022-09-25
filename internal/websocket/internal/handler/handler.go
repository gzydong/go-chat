package handler

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
	"go-chat/config"
	"go-chat/internal/pkg/im/adapter"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/pkg/logger"
)

type Handler struct {
	Chat    *ChatChannel
	Example *ExampleChannel
	Config  *config.Config
}

type AuthConn struct {
	Uid  int
	conn *adapter.TcpAdapter
}

type Authorize struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
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

func (h *Handler) auth(connect net.Conn, data chan *AuthConn) {
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

	if !gjson.GetBytes(read, "token").Exists() {
		return
	}

	detail := &Authorize{}
	if err := jsonutil.DecodeBt(read, detail); err != nil {
		return
	}

	claims, err := jwt.ParseToken(detail.Token, h.Config.Jwt.Secret)
	if err != nil || claims.Valid() != nil {
		return
	}

	uid, err := strconv.Atoi(claims.ID)
	if err != nil {
		return
	}

	if claims.Guard == "api" && detail.Channel == "chat" {
		data <- &AuthConn{Uid: uid, conn: conn}
	}
}
