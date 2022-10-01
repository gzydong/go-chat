package event

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"go-chat/internal/im_server/internal/event/example"
	"go-chat/internal/pkg/im"
)

type ExampleEvent struct {
	handler *example.Handler
}

func NewExampleEvent() *ExampleEvent {
	return &ExampleEvent{}
}

func (e *ExampleEvent) OnOpen(client im.IClient) {
	fmt.Printf("客户端[%d] 已连接\n", client.Cid())
}

func (e *ExampleEvent) OnMessage(client im.IClient, message []byte) {

	fmt.Println("接收消息===>>>", message)

	event := gjson.GetBytes(message, "event").String()
	if event != "" {
		// 触发事件
		e.handler.Call(context.Background(), client, event, message)
	}
}

func (e *ExampleEvent) OnClose(client im.IClient, code int, text string) {
	fmt.Printf("客户端[%d] 已关闭连接，关闭提示【%d】%s \n", client.Cid(), code, text)
}
