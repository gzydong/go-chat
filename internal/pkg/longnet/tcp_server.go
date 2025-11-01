// tcp_server.go

package longnet

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/gzydong/go-chat/internal/pkg/longnet/adapter"
)

type TcpServer struct {
	serv *Server
}

func newTcpServer(serv *Server) *TcpServer {
	return &TcpServer{serv}
}

func (t *TcpServer) Start(ctx context.Context) error {
	if t.serv.options.TCPConfig == nil {
		panic("tcp config is nil")
	}

	var err error
	var listener net.Listener
	if t.serv.options.TCPConfig.TLSEnable {
		listener, err = tls.Listen("tcp", t.serv.options.TCPConfig.Addr, t.serv.options.TLSConfig)
	} else {
		listener, err = net.Listen("tcp", t.serv.options.TCPConfig.Addr)
	}

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = listener.Close()
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				slog.Error("[tpc] accept error", "error", err)
				return
			}

			// 这里需要判断最大连接数，如果超出则返回错误
			if !t.serv.SessionManager().AllowAcceptConn() {
				_ = conn.Close()
				log.Printf("tcp connect error: %s", "too many connections")
				continue
			}

			// 这里先待定
			go t.handleConnection(conn)
		}
	}()

	slog.Info(fmt.Sprintf("Starting TCP server on %s", t.serv.options.TCPConfig.Addr))
	<-ctx.Done()
	slog.Info(fmt.Sprintf("TCP server on %s is shutting down...", t.serv.options.TCPConfig.Addr))
	return nil
}

type AuthorizeInfo struct {
	Event   string `json:"event"`
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func (t *TcpServer) handleConnection(conn net.Conn) {
	c, err := adapter.NewTcpAdapter(conn)
	if err != nil {
		slog.Error("Failed to create TCP adapter", "err", err)
		return
	}

	// 无需验证授权信息
	if t.serv.authorize == nil {
		t.serv.SessionManager().NewSession(0, c)
		return
	}

	// 设置读超时
	_ = c.SetReadDeadline(time.Now().Add(3 * time.Second))
	data, err := c.Read()
	if err != nil {
		_ = c.Close()
		slog.Error("Failed to read authorization info", "err", err)
		return
	}

	var authorizeInfo AuthorizeInfo
	if err := json.Unmarshal(data, &authorizeInfo); err != nil {
		_ = c.Write([]byte(`{"event":"unauthorized"}`))
		_ = c.Close()
		slog.Error("unmarshal authorize info err", "err", err)
		return
	}

	uid, err := t.serv.authorize(context.Background(), authorizeInfo.Payload.Token)
	if err != nil {
		_ = c.Write([]byte(`{"event":"unauthorized"}`))
		_ = c.Close()
		return
	}

	if err = c.Write([]byte(`{"event":"authorize"}`)); err != nil {
		_ = c.Close()
		slog.Error("write authorization response failed", "err", err)
		return
	}

	t.serv.SessionManager().NewSession(uid, c)
}
