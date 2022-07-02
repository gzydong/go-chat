package server

import (
	"context"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
)

var once sync.Once

type IServer interface {
	Setup(ctx context.Context) error
}

// SubServer 订阅的服务列表
type SubServer struct {
	Health    *Health      // 注册健康上报
	Subscribe *WsSubscribe // 注册消息订阅
}

type Server struct {
	items []IServer
}

func NewServer(routines *SubServer) *Server {
	server := &Server{}

	server.binds(routines)

	return server
}

func (c *Server) binds(routines *SubServer) {
	elem := reflect.ValueOf(routines).Elem()
	for i := 0; i < elem.NumField(); i++ {
		if v, ok := elem.Field(i).Interface().(IServer); ok {
			c.items = append(c.items, v)
		}
	}
}

// Start 启动服务
func (c *Server) Start(eg *errgroup.Group, ctx context.Context) {
	once.Do(func() {
		for _, process := range c.items {
			func(obj IServer) {
				eg.Go(func() error {
					return obj.Setup(ctx)
				})
			}(process)
		}
	})
}
