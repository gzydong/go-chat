package process

import (
	"context"
	"go-chat/app/pkg/im"
	"golang.org/x/sync/errgroup"
	"sync"
)

var onceProcess sync.Once

type InterfaceProcess interface {
	Handle(ctx context.Context) error
}

type Process struct {
	registers []InterfaceProcess
}

func NewProcessManage(serv *Server, subscribe *WsSubscribe, heart *Heartbeat, garbage *ClearGarbage) *Process {
	pro := &Process{}

	pro.Register(serv)
	pro.Register(subscribe)
	pro.Register(im.Sessions.Default)
	pro.Register(heart)
	pro.Register(garbage)

	return pro
}

func (p *Process) Register(process InterfaceProcess) {
	p.registers = append(p.registers, process)
}

func (p *Process) Run(eg *errgroup.Group, ctx context.Context) {
	onceProcess.Do(func() {
		for _, process := range p.registers {
			func(obj InterfaceProcess) {
				eg.Go(func() error {
					return obj.Handle(ctx)
				})
			}(process)
		}
	})
}
