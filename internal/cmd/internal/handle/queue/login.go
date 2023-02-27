package queue

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nsqio/go-nsq"
)

type LoginHandle struct {
	rds *redis.Client
}

func NewLoginHandle(rds *redis.Client) *LoginHandle {
	return &LoginHandle{rds: rds}
}

func (e *LoginHandle) Handle(ctx context.Context) error {

	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second

	c, err := nsq.NewConsumer("test", "test_channel", config)
	if err != nil {
		fmt.Printf("create consumer failed, err:%v\n", err)
		return err
	}

	defer c.Stop()

	c.AddHandler(e)

	if err := c.ConnectToNSQD(""); err != nil {
		return err
	}

	ch := make(chan os.Signal)        // 定义一个信号的通道
	signal.Notify(ch, syscall.SIGINT) // 转发键盘中断信号到c
	<-ch                              // 阻塞

	return nil
}

// HandleMessage 是需要实现的处理消息的方法
func (e *LoginHandle) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("%s recv from %v, msg:%v\n", "LoginHandle", msg.NSQDAddress, string(msg.Body))
	return nil
}
