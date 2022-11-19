package process

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

// SubServers 订阅的服务列表
type SubServers struct {
	HealthSubscribe  *HealthSubscribe  // 注册健康上报
	MessageSubscribe *MessageSubscribe // 注册消息订阅
}

type Server struct {
	items []IServer
}

func NewServer(servers *SubServers) *Server {
	s := &Server{}

	s.binds(servers)

	return s
}

func (c *Server) binds(servers *SubServers) {
	elem := reflect.ValueOf(servers).Elem()
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
			func(serv IServer) {
				eg.Go(func() error {
					return serv.Setup(ctx)
				})
			}(process)
		}
	})
}
