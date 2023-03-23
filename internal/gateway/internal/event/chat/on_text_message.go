package chat

import (
	"context"
	"encoding/json"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat/socket"
)

type TextMessage struct {
	MsgId   string                     `json:"msg_id"`
	AckId   string                     `json:"ack_id"`
	Event   string                     `json:"event"`
	Content message.TextMessageRequest `json:"content"`
}

// OnTextMessage 文本消息
func (h *Handler) OnTextMessage(ctx context.Context, client socket.IClient, data []byte) {

	var m TextMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnTextMessage Err: ", err)
		return
	}

	if m.Content.GetContent() == "" {
		return
	}

	if m.Content.GetReceiver() == nil {
		return
	}

	err := h.message.SendText(ctx, client.Uid(), &message.TextMessageRequest{
		Content: m.Content.Content,
		Receiver: &message.MessageReceiver{
			TalkType:   m.Content.Receiver.TalkType,
			ReceiverId: m.Content.Receiver.ReceiverId,
		},
	})

	if err != nil {
		log.Printf("Chat OnTextMessage err: %s", err.Error())
		return
	}

	err = client.Write(&socket.ClientResponse{
		Sid:   m.AckId,
		Event: "ack",
	})

	if err != nil {
		log.Printf("Chat OnTextMessage ack err: %s", err.Error())
		return
	}
}
