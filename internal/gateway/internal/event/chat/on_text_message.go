package chat

import (
	"context"
	"encoding/json"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat/socket"
)

type TextMessage struct {
	MsgId string                     `json:"msg_id"`
	Event string                     `json:"event"`
	Body  message.TextMessageRequest `json:"body"`
}

// OnTextMessage 文本消息
func (h *Handler) OnTextMessage(ctx context.Context, client socket.IClient, data []byte) {

	var m *TextMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnTextMessage Err: ", err)
		return
	}

	// fmt.Println("[TextMessage] 新消息 ", string(data))

	h.message.SendText(ctx, client.Uid(), &message.TextMessageRequest{
		Content: m.Body.Content,
		Receiver: &message.MessageReceiver{
			TalkType:   m.Body.Receiver.TalkType,
			ReceiverId: m.Body.Receiver.ReceiverId,
		},
	})

	// _ = client.Write(&im.ClientOutContent{
	// 	AckId: strutil.NewMsgId(),
	// 	IsAck: false,
	// 	Retry: 0,
	// 	Content: []byte(jsonutil.Encode(map[string]any{
	// 		"event":  "ack",
	// 		"ack_id": m.MsgId,
	// 	})),
	// })
}
