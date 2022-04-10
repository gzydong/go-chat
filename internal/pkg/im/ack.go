package im

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var ackManage = &AckBuffer{list: &sync.Map{}}

type AckBufferOption struct {
	Client  *Client // 客户端连接
	MsgID   string  // 消息ID
	Retry   int     // 重试次数
	Content []byte  // 内容
}

func AckManage() *AckBuffer {
	return ackManage
}

// AckBuffer Ack 确认缓冲区
type AckBuffer struct {
	list *sync.Map
}

func (a *AckBuffer) Add(opt *AckBufferOption) {
	a.list.Store(fmt.Sprintf("%s-%d", opt.MsgID, opt.Client.ClientId()), opt)
}

func (a *AckBuffer) Del(opt *AckBufferOption) {
	a.list.Delete(fmt.Sprintf("%s-%d", opt.MsgID, opt.Client.ClientId()))
}

func (a *AckBuffer) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		default:

			a.list.Range(func(key, value interface{}) bool {
				if option, ok := value.(*AckBufferOption); ok {
					_ = option.Client.Write(&ClientOutContent{
						IsAck:   true,
						Retry:   option.Retry,
						Content: option.Content,
					})

					a.Del(option)
				}

				return true
			})

			time.Sleep(30 * time.Second)
		}
	}
}
