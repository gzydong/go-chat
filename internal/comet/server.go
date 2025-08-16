package comet

import (
	"context"

	"go-chat/config"
	"go-chat/internal/pkg/longnet"
	"go-chat/internal/provider"
)

type Server struct {
	Config    *config.Config
	Subscribe *Subscribe
	Handler   *Handler
	Heartbeat *Heartbeat
	Authorize provider.UserJwtAuthorize
}

func (s *Server) Start(ctx context.Context) error {
	serv := longnet.New(longnet.Options{
		MaxOpenConns:  1000,
		MaxPacketSize: 2 << 20,
		WSSConfig: &longnet.WSSConfig{
			Addr: s.Config.Server.WebsocketAddr,
			Path: "/wss/default.io",
		},
	})
	serv.SetAuthorize(s.onAuthorize)
	serv.SetHandler(s.Handler)
	serv.SetEncoder(longnet.NewEncoder(longnet.EncoderOptions{
		MaxPacketSize:         2 << 20,   // 2M
		MinCompressPacketSize: 10 * 1024, // 10KB
	}, nil, nil))

	serv.SetCustomProcess(s.Heartbeat)
	serv.SetCustomProcess(s.Subscribe)

	return serv.Start(ctx)
}

// onTcpAuthorize 授权认证
func (s *Server) onAuthorize(ctx context.Context, token string) (int64, error) {
	claims, err := s.Authorize.Valid(token)
	if err != nil {
		return 0, err
	}

	return int64(claims.Metadata.UserId), nil
}
