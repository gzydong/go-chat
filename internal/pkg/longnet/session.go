package longnet

import (
	"sync"
	"sync/atomic"
	"time"
)

var _ ISession = (*Session)(nil)

type Session struct {
	mu           sync.Mutex
	connId       int64           // 会话ID
	userId       int64           // 用户ID
	lastActiveAt int64           // Unix 时间戳，单位为秒
	conn         IConn           // 连接
	closed       atomic.Bool     // 是否已关闭
	handler      IHandler        // 处理器
	manager      *SessionManager // 会话管理器
}

// NewSession 创建会话
func NewSession(uid int64, conn IConn, handler IHandler, manager *SessionManager) *Session {
	s := &Session{
		userId:       uid,
		connId:       manager.GenConnId(),
		conn:         conn,
		lastActiveAt: time.Now().Unix(),
		handler:      handler,
		manager:      manager,
	}

	return s
}

func (s *Session) init() {
	s.manager.Insert(s)

	go s.loopAccept()

	s.conn.SetCloseHandler(func(code int, text string) error {
		_ = s.Close()
		return nil
	})

	s.handler.OnOpen(s.manager, s)
}

func (s *Session) ConnId() int64 {
	return s.connId
}

func (s *Session) UserId() int64 {
	return s.userId
}

func (s *Session) IsClosed() bool {
	return s.closed.Load()
}

func (s *Session) Read() ([]byte, error) {
	_ = s.conn.SetReadDeadline(time.Now().Add(s.manager.Options().ReadTimeout))
	return s.conn.Read()
}

func (s *Session) Write(data []byte) (err error) {
	if s.IsClosed() {
		return ErrSessionClosed
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.conn.SetWriteDeadline(time.Now().Add(s.manager.Options().WriteTimeout))
	return s.conn.Write(data)
}

func (s *Session) Close() error {
	if s.closed.Swap(true) {
		return nil
	}

	_ = s.conn.Close()

	s.manager.Delete(s)

	s.handler.OnClose(s.connId, s.userId)
	return nil
}

func (s *Session) Network() string {
	return s.conn.Network()
}

func (s *Session) RefreshLastActiveAt() {
	s.lastActiveAt = time.Now().Unix()
}

func (s *Session) LastActiveAt() int64 {
	return s.lastActiveAt
}

// 循环接收客户端推送信息
func (s *Session) loopAccept() {
	defer func() {
		_ = s.Close()
	}()

	for {
		data, err := s.Read()
		if err != nil {
			break
		}

		s.RefreshLastActiveAt()

		s.handler.OnMessage(s.manager, s, data)
	}
}
