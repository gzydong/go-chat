package event

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"go-chat/internal/comet/handler/event/example"

	"go-chat/internal/pkg/core/socket"
)

type ExampleEvent struct {
	Handler *example.Handler
}

func (e *ExampleEvent) OnOpen(client socket.IClient) {
	fmt.Printf("客户端[%d] 已连接\n", client.Cid())
}

func (e *ExampleEvent) OnMessage(client socket.IClient, message []byte) {

	fmt.Println("接收消息===>>>", message)

	res := gjson.GetBytes(message, "event")
	if !res.Exists() {
		return
	}

	event := res.String()
	if event != "" {
		// 触发事件
		e.Handler.Call(context.TODO(), client, event, message)
	}
}

func (e *ExampleEvent) OnClose(client socket.IClient, code int, text string) {
	fmt.Printf("客户端[%d] 已关闭连接，关闭提示【%d】%s \n", client.Cid(), code, text)
}
