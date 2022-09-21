package chat

import "fmt"

// OnReadMessage 消息已读事件
func (h Handler) OnReadMessage(data string) {
	fmt.Println("OnReadMessage===>>>", data)
}
