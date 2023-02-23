package logic

import "context"

type Publisher struct {
}

func (p *Publisher) Publish(ctx context.Context, topic string, message any) error {
	return nil
}
