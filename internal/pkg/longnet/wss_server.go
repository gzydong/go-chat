package longnet

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gzydong/go-chat/internal/pkg/longnet/adapter"
)

type WssServer struct {
	serv *Server
}

func newWssServer(serv *Server) *WssServer {
	return &WssServer{
		serv: serv,
	}
}

func (s *WssServer) Start(ctx context.Context) error {
	options := s.serv.options

	mu := http.NewServeMux()

	mu.HandleFunc("GET "+options.WSSConfig.getPath(), func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		// 这里需要判断最大连接数，如果超出则返回错误
		if !s.serv.SessionManager().AllowAcceptConn() {
			w.WriteHeader(http.StatusTooManyRequests)
			log.Printf("[%s] websocket connect error: %s", r.RemoteAddr, "too many connections")
			return
		}

		var uid int64
		var err error
		if s.serv.authorize != nil {
			if uid, err = s.serv.authorize(context.Background(), token); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("[%s] websocket connect error: %s", r.RemoteAddr, err.Error())
				return
			}
		}

		conn, err := adapter.NewWsAdapter(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("[%s] websocket connect error: %s", r.RemoteAddr, err.Error())
			return
		}

		s.serv.SessionManager().NewSession(uid, conn)
	})

	server := http.Server{
		Addr:    options.WSSConfig.Addr,
		Handler: mu,
	}

	if options.WSSConfig.TLSEnable && options.TLSConfig != nil {
		server.TLSConfig = options.TLSConfig
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start listening: %s\n", err)
		}
	}()

	slog.Info(fmt.Sprintf("Starting WebSocket server on %s", options.WSSConfig.Addr))
	<-ctx.Done()
	slog.Info(fmt.Sprintf("WebSocket server on %s is shutting down...", options.WSSConfig.Addr))
	return server.Shutdown(ctx)
}
