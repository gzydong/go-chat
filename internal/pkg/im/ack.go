package im

import (
	"context"
	"sync"
	"time"
)

// nolint
var ack *AckBuffer

// AckBuffer Ack 确认缓冲区
type AckBuffer struct {
	node *Node
}

func init() {
	ack = &AckBuffer{NewNode(100)}
}

type AckBufferOption struct {
	Channel *Channel
	Cid     int64  // 客户端ID
	AckID   string // ACK ID
	Retry   int    // 重试次数
	Content []byte // 内容
}

// nolint
func (a *AckBuffer) add(opt *AckBufferOption) {
	a.node.nodes[a.node.index(opt.Cid)].Store(opt.AckID, opt)
}

// nolint
func (a *AckBuffer) del(opt *AckBufferOption) {
	a.node.nodes[a.node.index(opt.Cid)].Delete(opt.AckID)
}

func (a *AckBuffer) Start(ctx context.Context) error {

	timer := time.NewTimer(30 * time.Second)

	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:

			a.handle()

			timer.Reset(30 * time.Second)
		}
	}
}

func (a *AckBuffer) handle() {
	var sw sync.WaitGroup

	sw.Add(a.node.len)

	for _, v := range a.node.nodes {
		node := v

		go func() {
			defer sw.Done()

			node.Range(func(key, value any) bool {

				data := value.(AckBufferOption)

				client, isOk := data.Channel.Client(data.Cid)
				if !isOk {
					return true
				}

				_ = client.Write(&ClientOutContent{
					IsAck:   true,
					Retry:   data.Retry,
					Content: data.Content,
				})

				node.Delete(key)

				return true
			})
		}()
	}

	sw.Wait()
}
