package im

import "context"

type AckBufferOption struct {
	Client  Client // 客户端连接
	Content []byte // 内容
}

// nolint AckBuffer Ack 确认缓冲区
type AckBuffer struct {
	list []*AckBufferOption
}

func (a *AckBuffer) Add(opt *AckBufferOption) error {
	panic("implement me")
}

func (a *AckBuffer) Del(opt *AckBufferOption) error {
	panic("implement me")
}

func (a *AckBuffer) Run(ctx context.Context) error {
	panic("implement me")
}
