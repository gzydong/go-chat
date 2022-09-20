package consume

import (
	"github.com/sirupsen/logrus"
)

type ExampleSubscribe struct {
	handlers map[string]onConsumeFunc
}

func NewExampleSubscribe() *ExampleSubscribe {
	return &ExampleSubscribe{}
}

// Events 注册事件
func (s *ExampleSubscribe) Events() {
	s.handlers = make(map[string]onConsumeFunc)
}

// Call 触发回调事件
func (s *ExampleSubscribe) Call(event string, data string) {

	if s.handlers == nil {
		panic("事件未注册")
	}

	if f, ok := s.handlers[event]; ok {
		f(data)
	} else {
		logrus.Warnf("Event: [%s]未注册回调方法\n", event)
	}
}
