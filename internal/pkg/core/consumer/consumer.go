package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

type Consumer struct {
	addr      string
	config    *nsq.Config
	consumers map[string][]IConsumerHandle
}

type IConsumerHandle interface {
	Topic() string
	Channel() string
	Touch() bool
	Do(ctx context.Context, message []byte, attempts uint16) error
}

var deferred = map[int]int{
	1: 5,
	2: 15,
	3: 30,
	4: 180,
	5: 1800,
	6: 1800,
	7: 1800,
	8: 1800,
	9: 3600,
}

type BackoffStrategy struct {
}

func (b *BackoffStrategy) Calculate(attempt int) time.Duration {
	delay, ok := deferred[attempt]
	if !ok {
		return -1
	}

	return time.Duration(delay) * time.Second
}

func NewConsumer(addr string, conf *nsq.Config) *Consumer {
	conf.MaxInFlight = 20                     // 最大并发处理的消息数
	conf.HeartbeatInterval = 15 * time.Second // 心跳间隔
	conf.ReadTimeout = 20 * time.Second       // 读取超时
	conf.WriteTimeout = 20 * time.Second      // 写入超时

	return &Consumer{
		addr:      addr,
		config:    conf,
		consumers: make(map[string][]IConsumerHandle),
	}
}

func (c *Consumer) Register(group string, handle IConsumerHandle) {
	if _, ok := c.consumers[group]; !ok {
		c.consumers[group] = make([]IConsumerHandle, 0)
	}

	c.consumers[group] = append(c.consumers[group], handle)
}

func (c *Consumer) Start(ctx context.Context, group string) error {
	items, ok := c.consumers[group]
	if !ok {
		return fmt.Errorf("consumer group [%s] not found", group)
	}

	for _, item := range items {
		go c.start(ctx, item)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
	case <-ctx.Done():
	}

	// 等待5秒钟
	time.Sleep(time.Second * 5)

	return nil
}

func (c *Consumer) start(ctx context.Context, handle IConsumerHandle) {
	consumer, err := nsq.NewConsumer(handle.Topic(), handle.Channel(), c.config)
	if err != nil {
		panic(fmt.Errorf("[Consumer] NewConsumer error: %v", err))
	}

	strategy := &BackoffStrategy{}
	consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(message *nsq.Message) error {
		message.DisableAutoResponse()

		if handle.Touch() {
			timer := time.NewTimer(time.Second * 10)
			go func() {
				<-timer.C
				message.Touch()
			}()

			defer timer.Stop()
		}

		err = handle.Do(context.Background(), message.Body, message.Attempts)
		if err == nil {
			message.Finish()
			return nil
		}

		message.RequeueWithoutBackoff(strategy.Calculate(int(message.Attempts)))
		return nil
	}), 100)

	if err = consumer.ConnectToNSQD(c.addr); err != nil {
		panic(fmt.Errorf("[Consumer] ConnectToNSQD error: %v", err))
	}

	// 等待退出信号
	<-ctx.Done()

	// 停止消费
	consumer.Stop()

	// 阻塞等待消息处理
	<-consumer.StopChan
}
