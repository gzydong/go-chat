package longnet

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gzydong/go-chat/internal/pkg/server"
	"golang.org/x/sync/errgroup"
)

var _ IServer = (*Server)(nil)

type Server struct {
	serverId    string // 服务ID
	ctx         context.Context
	cancel      context.CancelFunc
	eg          errgroup.Group
	options     *Options
	process     []IProcess
	authorize   IAuthorize
	handler     IHandler
	idGenerator IdGenerator
	encoder     IEncoder
	manager     *SessionManager
}

func New(c Options) *Server {
	c = c.init()

	s := &Server{
		serverId: server.ID(),
		options:  &c,
	}

	return s
}

func (s *Server) ServerId() string {
	return s.serverId
}

func (s *Server) SetCustomProcess(process IProcess) {
	s.process = append(s.process, process)
}

func (s *Server) SetAuthorize(auth IAuthorize) {
	s.authorize = auth
}

func (s *Server) SetHandler(h IHandler) {
	s.handler = h
}

func (s *Server) SessionManager() ISessionManager {
	return s.manager
}

func (s *Server) SetEncoder(encoder IEncoder) {
	s.encoder = encoder
}

func (s *Server) Encoder() IEncoder {
	return s.encoder
}

func (s *Server) Handler() IHandler {
	return s.handler
}

func (s *Server) SetIdGenerator(gen IdGenerator) {
	s.idGenerator = gen
}

func (s *Server) IdGenerator() IdGenerator {
	return s.idGenerator
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	defer s.cancel()
	defer slog.Info("Server shutdown completed.")

	slog.Info("Server is starting up...")
	slog.Info(fmt.Sprintf("Server ID: %s", s.serverId))
	slog.Info(fmt.Sprintf("Process ID(PID): %d", os.Getpid()))

	s.init()
	s.initWssServer()
	s.initTcpServer()
	s.initProcess()

	// 监听关闭信号
	s.listenShutdown()

	return s.eg.Wait()
}

func (s *Server) init() {
	s.manager = newSessionManager(s.options, s)

	if s.idGenerator == nil {
		s.idGenerator = NewAutoIdGenerator()
	}

	if s.encoder == nil {
		s.encoder = NewEncoder(EncoderOptions{}, nil, nil)
	}

	_ = s.manager.Start(s.ctx) // 启动会话管理器
}

func (s *Server) initWssServer() {
	if s.options.WSSConfig != nil {
		s.eg.Go(func() error {
			ws := newWssServer(s)
			return ws.Start(s.ctx)
		})
	}
}

func (s *Server) initTcpServer() {
	if s.options.TCPConfig != nil {
		s.eg.Go(func() error {
			tcp := newTcpServer(s)
			return tcp.Start(s.ctx)
		})
	}
}

func (s *Server) initProcess() {
	for _, v := range s.process {
		s.eg.Go(func() error {
			return v.Start(s.ctx, s)
		})
	}
}

func (s *Server) listenShutdown() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		<-c
		log.Println("Shutting down server...")
		s.cancel()
	}()
}
