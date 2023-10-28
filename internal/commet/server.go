package commet

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/commet/handler"
	"go-chat/internal/commet/process"
	"go-chat/internal/pkg/email"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/provider"
	"golang.org/x/sync/errgroup"
)

var ErrServerClosed = errors.New("shutting down server")

type AppProvider struct {
	Config    *config.Config
	Engine    *gin.Engine
	Coroutine *process.Server
	Handler   *handler.Handler
	Providers *provider.Providers
}

func Run(ctx *cli.Context, app *AppProvider) error {
	eg, groupCtx := errgroup.WithContext(ctx.Context)

	if !app.Config.Debug() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化 IM 渠道配置
	socket.Initialize(groupCtx, eg, func(name string) {
		emailClient := app.Providers.EmailClient
		if app.Config.App.Env == "prod" {
			_ = emailClient.SendMail(&email.Option{
				To:      app.Config.App.AdminEmail,
				Subject: fmt.Sprintf("[%s]守护进程异常", app.Config.App.Env),
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

	log.Printf("Server ID   :%s", app.Config.ServerId())
	log.Printf("Server Pid  :%d", os.Getpid())
	log.Printf("Websocket Listen Port :%d", app.Config.Server.Websocket)
	log.Printf("Tcp Listen Port :%d", app.Config.Server.Tcp)
	go NewTcpServer(app)

	return start(c, eg, groupCtx, app)
}

func start(c chan os.Signal, eg *errgroup.Group, ctx context.Context, app *AppProvider) error {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.Websocket),
		Handler: app.Engine,
	}

	eg.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	eg.Go(func() (err error) {
		defer func() {
			log.Println("Shutting down server...")

			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.TODO(), 3*time.Second)
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

	log.Println("Server exiting")

	return nil
}

func NewTcpServer(app *AppProvider) {

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", app.Config.Server.Tcp))
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept failed err:", err)
			continue
		}

		// TCP 分发
		go app.Handler.Dispatch(conn)
	}
}
