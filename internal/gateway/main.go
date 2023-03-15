package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/pkg/email"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	cmd := cli.NewApp()

	cmd.Name = "LumenIM 在线聊天"
	cmd.Usage = "IM Server"

	// 设置参数
	cmd.Flags = []cli.Flag{
		// 配置文件参数
		&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.yaml", Usage: "配置文件路径", DefaultText: "./config.yaml"},
	}

	cmd.Action = newApp

	_ = cmd.Run(os.Args)
}

func newApp(tx *cli.Context) error {
	eg, groupCtx := errgroup.WithContext(tx.Context)

	// 读取配置文件
	conf := config.ReadConfig(tx.String("config"))

	// 设置日志输出
	logger.SetOutput(conf.GetLogPath(), "logger-ws")

	if !conf.Debug() {
		gin.SetMode(gin.ReleaseMode)
	}

	app := Initialize(conf)

	// 初始化 IM 渠道配置
	socket.Initialize(groupCtx, eg, func(name string) {
		emailClient := app.Providers.EmailClient
		if conf.App.Env == "prod" {
			_ = emailClient.SendMail(&email.Option{
				To:      []string{"837215079@qq.com"},
				Subject: fmt.Sprintf("[%s]守护进程异常", conf.App.Env),
				Body:    fmt.Sprintf("守护进程异常[%s]", name),
			})
		}
	})

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// 延时启动守护协程
	time.AfterFunc(3*time.Second, func() {
		app.Coroutine.Start(eg, groupCtx)
	})

	log.Printf("Server ID   :%s", conf.ServerId())
	log.Printf("Server Pid  :%d", os.Getpid())
	log.Printf("Websocket Listen Port :%d", conf.Server.Websocket)
	log.Printf("Tcp Listen Port :%d", conf.Server.Tcp)

	go NewTcpServer(app)

	return start(c, eg, groupCtx, app)
}

var ErrServerClosed = errors.New("shutting down server")

func start(c chan os.Signal, eg *errgroup.Group, ctx context.Context, app *AppProvider) error {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.Websocket),
		Handler: app.Engine,
	}

	eg.Go(func() error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	eg.Go(func() (err error) {
		defer func() {
			log.Println("Shutting down server...")

			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer timeCancel()

			if err := server.Shutdown(timeCtx); err != nil {
				log.Printf("Websocket Shutdown Err: %s \n", err)
			}

			err = ErrServerClosed
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return nil
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, ErrServerClosed) {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	time.Sleep(3 * time.Second)
	log.Println("Server exiting")

	return nil
}
