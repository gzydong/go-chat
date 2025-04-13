package adapter

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WsAdapter Websocket 适配器
type WsAdapter struct {
	conn          *websocket.Conn
	readDeadline  time.Duration
	writeDeadline time.Duration
}

var defaultUpGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWsAdapter(w http.ResponseWriter, r *http.Request) (*WsAdapter, error) {
	conn, err := defaultUpGrader.Upgrade(w, r, w.Header())
	if err != nil {
		return nil, err
	}

	return &WsAdapter{
		conn:          conn,
		readDeadline:  3 * time.Minute,
		writeDeadline: 3 * time.Second,
	}, nil
}

func (w *WsAdapter) Network() string {
	return NetworkWss
}

func (w *WsAdapter) Read() ([]byte, error) {
	if w.readDeadline > 0 {
		_ = w.conn.SetReadDeadline(time.Now().Add(w.readDeadline))
	}

	_, content, err := w.conn.ReadMessage()
	return content, err
}

func (w *WsAdapter) Write(bytes []byte) error {
	if w.writeDeadline > 0 {
		_ = w.conn.SetWriteDeadline(time.Now().Add(w.writeDeadline))
	}

	return w.conn.WriteMessage(websocket.TextMessage, bytes)
}

func (w *WsAdapter) Close() error {
	return w.conn.Close()
}

func (w *WsAdapter) SetCloseHandler(fn func(code int, text string) error) {
	w.conn.SetCloseHandler(fn)
}

func (w *WsAdapter) SetReadDeadline(duration time.Duration) {
	w.readDeadline = duration
}

func (w *WsAdapter) SetWriteDeadline(duration time.Duration) {
	w.writeDeadline = duration
}
