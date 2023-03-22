package socket

import (
	"log"

	"go-chat/internal/pkg/utils"
)

type IEvent interface {
	Open(client IClient)
	Message(client IClient, data []byte)
	Close(client IClient, code int, text string)
	Destroy(client IClient)
}

type (
	OpenEvent         func(client IClient)
	MessageEvent      func(client IClient, data []byte)
	CloseEvent        func(client IClient, code int, text string)
	DestroyEvent      func(client IClient)
	ClientEventOption func(event *ClientEvent)
)

type ClientEvent struct {
	open    OpenEvent
	message MessageEvent
	close   CloseEvent
	destroy DestroyEvent
}

func NewClientEvent(opts ...ClientEventOption) IEvent {

	o := &ClientEvent{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (c *ClientEvent) Open(client IClient) {

	if c.open == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("open event callback exception: ", client.Uid(), client.Cid(), client.ChannelName(), utils.PanicTrace(err))
		}
	}()

	c.open(client)
}

func (c *ClientEvent) Message(client IClient, data []byte) {

	if c.message == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("message event callback exception: ", client.Uid(), client.Cid(), client.ChannelName(), utils.PanicTrace(err))
		}
	}()

	c.message(client, data)
}

func (c *ClientEvent) Close(client IClient, code int, text string) {
	if c.close == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("close event callback exception: ", client.Uid(), client.Cid(), client.ChannelName(), utils.PanicTrace(err))
		}
	}()

	c.close(client, code, text)
}

func (c *ClientEvent) Destroy(client IClient) {
	if c.destroy == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("destroy event callback exception: ", client.Uid(), client.Cid(), client.ChannelName(), utils.PanicTrace(err))
		}
	}()

	c.destroy(client)
}

// WithOpenCallback 连接成功回调事件
func WithOpenCallback(e OpenEvent) ClientEventOption {
	return func(event *ClientEvent) {
		event.open = e
	}
}

// WithMessageCallback 消息回调事件
func WithMessageCallback(e MessageEvent) ClientEventOption {
	return func(event *ClientEvent) {
		event.message = e
	}
}

// WithCloseCallback 连接关闭回调事件
func WithCloseCallback(e CloseEvent) ClientEventOption {
	return func(event *ClientEvent) {
		event.close = e
	}
}

// WithDestroyCallback 连接销毁回调事件
func WithDestroyCallback(e DestroyEvent) ClientEventOption {
	return func(event *ClientEvent) {
		event.destroy = e
	}
}
