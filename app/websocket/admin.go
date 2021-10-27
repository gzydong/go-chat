package websocket

import (
	"fmt"
	"go-chat/app/pkg/im"
)

type AdminChannelHandle struct {
}

func NewAdminChannelHandle() *AdminChannelHandle {
	return &AdminChannelHandle{}
}

func (d *AdminChannelHandle) Open(client *im.Client) {

}

func (d *AdminChannelHandle) Message(message *im.RecvMessage) {
	fmt.Println("ws-socket", message.Content)
}

func (d *AdminChannelHandle) Close(client *im.Client, code int, text string) {

}
