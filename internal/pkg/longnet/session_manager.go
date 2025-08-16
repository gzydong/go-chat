package longnet

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var _ ISessionManager = (*SessionManager)(nil)

type SessionManagerOption func(sm *SessionManager)

type SessionManager struct {
	options     *Options                            // 会话配置
	currConnNum AtomicInt32                         // 当前连接数
	heartbeat   IHeartbeat                          // 心跳管理器
	sessions    cmap.ConcurrentMap[int64, ISession] // 会话管理
	userSession *SetShards                          // 用户会话管理
	assistant   IServerAssist                       // 辅助
}

func newSessionManager(options *Options, assistant IServerAssist) *SessionManager {
	session := &SessionManager{
		options:     options,
		currConnNum: AtomicInt32{},
		sessions:    cmap.NewWithCustomShardingFunction[int64, ISession](fnv32),
		userSession: NewSetShards(),
		assistant:   assistant,
	}

	return session
}

func (s *SessionManager) AllowAcceptConn() bool {
	return !(s.options.MaxOpenConns > 0 && int(s.GetSessionNum()) >= s.options.MaxOpenConns)
}

func (s *SessionManager) Options() *Options {
	return s.options
}

func (s *SessionManager) GenConnId() int64 {
	return s.assistant.IdGenerator().IdGen()
}

func (s *SessionManager) NewSession(uid int64, conn IConn) {
	NewSession(uid, conn, s.assistant.Handler(), s).init()
}

func (s *SessionManager) Insert(c ISession) {
	s.currConnNum.Incr()

	s.sessions.Set(c.ConnId(), c)

	s.heartbeat.Insert(c.ConnId(), s.options.PingInterval)
	if c.UserId() > 0 {
		s.userSession.Add(c.UserId(), c.ConnId())
	}
}

func (s *SessionManager) Delete(c ISession) {
	s.heartbeat.Cancel(c.ConnId())

	s.sessions.Remove(c.ConnId())
	s.currConnNum.Decr()

	if c.UserId() > 0 {
		s.userSession.Del(c.UserId(), c.ConnId())
	}
}

func (s *SessionManager) GetSession(cid int64) (ISession, error) {
	session, ok := s.sessions.Get(cid)
	if !ok {
		return nil, ErrSessionNotExist
	}

	return session, nil
}

func (s *SessionManager) GetConnIds(uid int64) []int64 {
	return s.userSession.Get(uid)
}

func (s *SessionManager) GetSessions(uid int64) []ISession {
	items := make([]ISession, 0)
	for _, cid := range s.GetConnIds(uid) {
		session, ok := s.sessions.Get(cid)
		if !ok {
			continue
		}

		items = append(items, session)
	}

	return items
}

func (s *SessionManager) GetSessionNum() int32 {
	return s.currConnNum.Load()
}

func (s *SessionManager) GetSessionUserNum() int32 {
	return s.userSession.GetUserNum()
}

func (s *SessionManager) Assistant() IServerAssist {
	return s.assistant
}

func (s *SessionManager) Iterator() <-chan ISession {
	ch := make(chan ISession)

	go func() {
		defer close(ch)

		for session := range s.sessions.IterBuffered() {
			ch <- session.Val
		}
	}()

	return ch
}

func (s *SessionManager) onHeartbeatCallback(cid int64) {
	session, ok := s.sessions.Get(cid)
	if !ok {
		slog.Warn(fmt.Sprintf("session not exist: %d", cid))
		return
	}

	if session.IsClosed() {
		slog.Warn(fmt.Sprintf("session closed: %d", cid))
		return
	}

	now := time.Now().Unix()
	lastActiveAt := session.LastActiveAt()

	// 心跳检测超时
	if now-lastActiveAt > int64(s.options.PingTimeout) {
		slog.Debug(fmt.Sprintf("session %d timeout, last active: %d, now: %d", cid, lastActiveAt, now))

		if err := session.Close(); err != nil {
			slog.Warn(fmt.Sprintf("session %d close error: %v", cid, err))
		}

		return
	}

	s.heartbeat.Insert(cid, s.options.PingInterval)

	if now-lastActiveAt > 60 {
		if err := session.Write([]byte(`{"cmd":"ping"}`)); err != nil {
			slog.Error(fmt.Sprintf("session %d ping error: %v", cid, err))
		}
	}
}

func (s *SessionManager) Start(ctx context.Context) error {
	s.heartbeat = NewHeartbeat(30, s.onHeartbeatCallback)
	return nil
}
