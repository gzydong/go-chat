package logic

import "context"

type Publisher struct {
}

func (p *Publisher) Publish(ctx context.Context, topic string, message interface{}) error {
	return nil
}
