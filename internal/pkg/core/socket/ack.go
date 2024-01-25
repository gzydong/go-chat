package socket

import (
	"context"
	"errors"
	"log"
	"time"

	"go-chat/internal/pkg/timewheel"
)

var ack *AckBuffer

// AckBuffer Ack 确认缓冲区
type AckBuffer struct {
	timeWheel *timewheel.SimpleTimeWheel[*AckBufferContent]
}

type AckBufferContent struct {
	cid      int64
	uid      int64
	channel  string
	response *ClientResponse
}

func InitAck() {
	ack = &AckBuffer{}
	ack.timeWheel = timewheel.NewSimpleTimeWheel[*AckBufferContent](1*time.Second, 30, ack.handle)
}

func (a *AckBuffer) Start(ctx context.Context) error {

	go a.timeWheel.Start()

	<-ctx.Done()

	a.timeWheel.Stop()

	return errors.New("ack service stopped")
}

func (a *AckBuffer) insert(ackKey string, value *AckBufferContent) {
	a.timeWheel.Add(ackKey, value, time.Duration(5)*time.Second)
}

func (a *AckBuffer) delete(ackKey string) {
	a.timeWheel.Remove(ackKey)
}

func (a *AckBuffer) handle(_ *timewheel.SimpleTimeWheel[*AckBufferContent], _ string, bufferContent *AckBufferContent) {

	ch, ok := Session.Channel(bufferContent.channel)
	if !ok {
		return
	}

	client, ok := ch.Client(bufferContent.cid)
	if !ok {
		return
	}

	if client.Closed() || int64(client.uid) != bufferContent.uid {
		return
	}

	if err := client.Write(bufferContent.response); err != nil {
		log.Println("ack err: ", err)
	}
}
