package socket

import (
	"context"
	"errors"
	"log"
	"time"

	"go-chat/internal/pkg/timewheel"
)

// nolint
var ack *AckBuffer

// AckBuffer Ack 确认缓冲区
type AckBuffer struct {
	timeWheel *timewheel.SimpleTimeWheel
}

type AckBufferBody struct {
	Cid   int64
	Uid   int64
	Ch    string
	Value *ClientResponse
}

func init() {
	ack = &AckBuffer{}
	ack.timeWheel = timewheel.NewSimpleTimeWheel(1*time.Second, 30, ack.handle)
}

func (a *AckBuffer) Start(ctx context.Context) error {

	go a.timeWheel.Start()

	<-ctx.Done()

	a.timeWheel.Stop()

	return errors.New("AckBuffer exit")
}

func (a *AckBuffer) add(ackKey string, value *AckBufferBody) {
	_ = a.timeWheel.Add(ackKey, value, time.Duration(5)*time.Second)
}

// nolint
func (a *AckBuffer) remove(ackKey string) {
	a.timeWheel.Remove(ackKey)
}

func (a *AckBuffer) handle(_ *timewheel.SimpleTimeWheel, key string, value any) {
	buffer, ok := value.(*AckBufferBody)
	if !ok {
		return
	}

	ch, ok := Session.Channel(buffer.Ch)
	if !ok {
		return
	}

	// 重发消息，需要检测客户端是否已断开，如果已断开则不需要重发
	client, ok := ch.Client(buffer.Cid)
	if !ok {
		return
	}

	if client.Closed() {
		return
	}

	if int64(client.uid) != buffer.Uid {
		return
	}

	err := client.Write(buffer.Value)
	if err != nil {
		log.Println("AckBuffer ack err: ", err)
	}
}
