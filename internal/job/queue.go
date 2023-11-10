package job

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/job/queue"
	"gorm.io/gorm"
)

type QueueProvider struct {
	Config *config.Config
	DB     *gorm.DB
	Queue  *Queue
}

type Queue struct {
	queue.ExampleQueue
}

func RunQueue(ctx *cli.Context, app *QueueProvider) error {
	log.Println("队列运行中...")

	err := app.Queue.ExampleQueue.Handle(ctx.Context)
	if err != nil {
		fmt.Println("ExampleQueue>>", err)
	}

	ch := make(chan os.Signal, 1)     // 定义一个信号的通道
	signal.Notify(ch, syscall.SIGINT) // 转发键盘中断信号到c
	<-ch                              // 阻塞

	log.Println("队列已结束...")

	return nil
}
