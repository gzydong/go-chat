package adapter

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WsAdapter Websocket 适配器
type WsAdapter struct {
	conn *websocket.Conn
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

	return &WsAdapter{conn: conn}, nil
}

func (w *WsAdapter) Network() string {
	return NetworkWss
}

func (w *WsAdapter) Read() ([]byte, error) {
	_, content, err := w.conn.ReadMessage()
	return content, err
}

func (w *WsAdapter) Write(bytes []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, bytes)
}

func (w *WsAdapter) Close() error {
	return w.conn.Close()
}

func (w *WsAdapter) SetCloseHandler(fn func(code int, text string) error) {
	w.conn.SetCloseHandler(fn)
}
