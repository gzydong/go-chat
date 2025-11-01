package longnet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gzydong/go-chat/internal/pkg/longnet/adapter/encoding"
	"github.com/pkg/errors"
)

var ErrClientClosed = errors.New("client closed")

type TcpClient struct {
	addr      string
	w         net.Conn
	r         *bufio.Reader
	closed    atomic.Bool
	mu        sync.Mutex
	compress  ICompress
	onMessage func(c IClient, data []byte)
}

type Command struct {
	Event   string `json:"event,omitempty"`
	Payload any    `json:"payload,omitempty"`
}

func NewTcpClient(addr string) *TcpClient {
	return &TcpClient{
		addr: addr,
	}
}

func (t *TcpClient) Connect(token string) error {
	conn, err := net.DialTimeout("tcp", t.addr, 3*time.Second)
	if err != nil {
		return ErrClientClosed
	}

	t.w = conn
	t.r = bufio.NewReader(conn)

	data, _ := encoding.NewEncode(NewAuthorizeCommand(token))
	if _, err = t.w.Write(data); err != nil {
		t.Close()
		slog.Error("authorize write error", "err", err)
		return err
	}

	resp, err := t.read()
	if err != nil {
		_ = conn.Close()
		return err
	}

	if !isAuthorize(resp) {
		t.Close()
		slog.Warn("authorization failed", "resp", string(resp))
		return errors.New("authorization failed")
	}

	go t.loopRead()
	return nil
}

func (t *TcpClient) Write(data []byte) error {
	if t.w == nil {
		return errors.New("client not connected")
	}

	if t.closed.Load() {
		return errors.New("client closed")
	}

	var err error
	if t.compress != nil {
		data, err = t.compress.Compress(data)
		if err != nil {
			return err
		}
	}

	value, err := encoding.NewEncode(data)
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	_, err = t.w.Write(value)
	return err
}

func (t *TcpClient) Close() {
	if t.closed.Swap(true) {
		return
	}

	_ = t.w.Close()
}

func (t *TcpClient) SetOnMessage(fn func(client IClient, data []byte)) {
	t.onMessage = fn
}

func (t *TcpClient) SetCompress(compress ICompress) {
	t.compress = compress
}

func (t *TcpClient) read() ([]byte, error) {
	msg, err := encoding.NewDecode(t.r)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (t *TcpClient) loopRead() {
	for {
		data, err := t.read()
		if err != nil {
			t.Close()
			return
		}

		if t.compress != nil {
			data, err = t.compress.Decompress(data)
			if err != nil {
				slog.Error("decompress err", "err", err)
				return
			}
		}

		if t.onMessage != nil {
			t.onMessage(t, data)
		}
	}
}

func NewAuthorizeCommand(token string) []byte {
	return []byte(fmt.Sprintf(`{"event":"%s","payload":{"token":"%s"}}`, token, token))
}

func isAuthorize(data []byte) bool {
	var info struct {
		Event string `json:"event"`
	}

	err := json.Unmarshal(data, &info)
	if err != nil {
		slog.Error("json.Unmarshal err", "err", err)
		return false
	}

	return info.Event == "authorize"
}
