package process

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

var onceProcess sync.Once

type InterfaceProcess interface {
	Handle(ctx context.Context) error
}

type Process struct {
	registers []InterfaceProcess
}

func NewProcess(health *Health, subscribe *WsSubscribe) *Process {
	process := &Process{}

	process.register(health)
	process.register(subscribe)

	return process
}

func (p *Process) register(process InterfaceProcess) {
	p.registers = append(p.registers, process)
}

func (p *Process) Start(eg *errgroup.Group, ctx context.Context) {
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
