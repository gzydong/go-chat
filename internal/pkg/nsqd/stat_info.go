package nsqd

import (
	"context"
	"fmt"
)

func GetStatInfo(ctx context.Context) (*NsqStatInfo, error) {
	return doGet[NsqStatInfo](ctx, "http://127.0.0.1:4151/stats?format=json&include_mem=false")
}

type deleteMessage struct {
	Message string `json:"message"`
}

func TopicDelete(ctx context.Context, topic string) error {
	res, err := doDelete[deleteMessage](ctx, fmt.Sprintf("http://127.0.0.1:4171/api/topics/%s", topic))
	if err != nil {
		return err
	}

	if res.Message != "" {
		return fmt.Errorf(res.Message)
	}

	return nil
}

func ChannelDelete(ctx context.Context, topic string, channel string) error {
	res, err := doDelete[deleteMessage](ctx, fmt.Sprintf("http://127.0.0.1:4171/api/topics/%s/%s", topic, channel))

	if err != nil {
		return err
	}

	if res.Message != "" {
		return fmt.Errorf(res.Message)
	}

	return nil
}
