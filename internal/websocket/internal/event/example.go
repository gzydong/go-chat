package event

import (
	"fmt"

	"go-chat/internal/pkg/im"
)

type ExampleEvent struct {
}

func NewExampleEvent() *ExampleEvent {
	return &ExampleEvent{}
}

func (e *ExampleEvent) OnOpen(client im.IClient) {
	fmt.Printf("客户端[%d] 已连接\n", client.Cid())
}

func (e *ExampleEvent) OnMessage(client im.IClient, message []byte) {
	fmt.Println("接收消息===>>>", message)
}

func (e *ExampleEvent) OnClose(client im.IClient, code int, text string) {
	fmt.Printf("客户端[%d] 已关闭连接，关闭提示【%d】%s \n", client.Cid(), code, text)
}
