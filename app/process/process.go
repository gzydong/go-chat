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

func NewProcessManage(run *ServerRun, ws *WsSubscribe) *Process {
	pro := &Process{}

	pro.Register(run)
	pro.Register(ws)
	pro.Register(im.Session.DefaultChannel)

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
