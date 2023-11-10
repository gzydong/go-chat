package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat/socket"
)

var publishMapping map[string]handle

func (h *Handler) onPublish(ctx context.Context, client socket.IClient, data []byte) {
	if publishMapping == nil {
		publishMapping = make(map[string]handle)
		publishMapping["text"] = h.onTextMessage
		publishMapping["code"] = h.onCodeMessage
		publishMapping["location"] = h.onLocationMessage
		publishMapping["emoticon"] = h.onEmoticonMessage
		publishMapping["vote"] = h.onVoteMessage
		publishMapping["image"] = h.onImageMessage
		publishMapping["file"] = h.onFileMessage
	}

	val, err := sonic.Get(data, "content.type")
	if err == nil {
		return
	}

	// 获取事件名
	typeValue, _ := val.String()
	if call, ok := publishMapping[typeValue]; ok {
		call(ctx, client, data)
	} else {
		log.Printf("chat event: onPublish [%s]未知的消息类型\n", typeValue)
	}
}

type TextMessage struct {
	AckId   string                     `json:"ack_id"`
	Event   string                     `json:"event"`
	Content message.TextMessageRequest `json:"content"`
}

// onTextMessage 文本消息
func (h *Handler) onTextMessage(ctx context.Context, client socket.IClient, data []byte) {

	var in TextMessage
	if err := json.Unmarshal(data, &in); err != nil {
		log.Println("Chat onTextMessage Err: ", err)
		return
	}

	if in.Content.GetContent() == "" || in.Content.GetReceiver() == nil {
		return
	}

	err := h.MessageService.SendText(ctx, client.Uid(), &message.TextMessageRequest{
		Content: in.Content.Content,
		Receiver: &message.MessageReceiver{
			TalkType:   in.Content.Receiver.TalkType,
			ReceiverId: in.Content.Receiver.ReceiverId,
		},
	})

	if err != nil {
		log.Printf("Chat onTextMessage err: %s", err.Error())
		return
	}

	if len(in.AckId) == 0 {
		return
	}

	if err = client.Write(&socket.ClientResponse{Sid: in.AckId, Event: "ack"}); err != nil {
		log.Printf("Chat onTextMessage ack err: %s", err.Error())
	}
}

type CodeMessage struct {
	AckId   string                     `json:"ack_id"`
	Event   string                     `json:"event"`
	Content message.CodeMessageRequest `json:"content"`
}

// onCodeMessage 代码消息
func (h *Handler) onCodeMessage(ctx context.Context, client socket.IClient, data []byte) {

	var m CodeMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat onTextMessage Err: ", err)
		return
	}

	if m.Content.GetReceiver() == nil {
		return
	}

	err := h.MessageService.SendCode(ctx, client.Uid(), &message.CodeMessageRequest{
		Lang: m.Content.Lang,
		Code: m.Content.Code,
		Receiver: &message.MessageReceiver{
			TalkType:   m.Content.Receiver.TalkType,
			ReceiverId: m.Content.Receiver.ReceiverId,
		},
	})

	if err != nil {
		log.Printf("Chat onTextMessage err: %s", err.Error())
		return
	}

	if len(m.AckId) == 0 {
		return
	}

	if err = client.Write(&socket.ClientResponse{Sid: m.AckId, Event: "ack"}); err != nil {
		log.Printf("Chat onTextMessage ack err: %s", err.Error())
	}
}

type EmoticonMessage struct {
	MsgId   string                         `json:"msg_id"`
	Event   string                         `json:"event"`
	Content message.EmoticonMessageRequest `json:"content"`
}

// onEmoticonMessage 表情包消息
func (h *Handler) onEmoticonMessage(_ context.Context, _ socket.IClient, data []byte) {

	var m EmoticonMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat onEmoticonMessage Err: ", err)
		return
	}

	fmt.Println("[onEmoticonMessage] 新消息 ", string(data))
}

type ImageMessage struct {
	MsgId   string                      `json:"msg_id"`
	Event   string                      `json:"event"`
	Content message.ImageMessageRequest `json:"content"`
}

// onImageMessage 图片消息
func (h *Handler) onImageMessage(_ context.Context, _ socket.IClient, data []byte) {

	var m ImageMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat onImageMessage Err: ", err)
		return
	}

	fmt.Println("[onImageMessage] 新消息 ", string(data))
}

type FileMessage struct {
	MsgId   string                      `json:"msg_id"`
	Event   string                      `json:"event"`
	Content message.ImageMessageRequest `json:"content"`
}

// onFileMessage 文件消息
func (h *Handler) onFileMessage(_ context.Context, _ socket.IClient, data []byte) {

	var m FileMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat onFileMessage Err: ", err)
		return
	}

	fmt.Println("[onFileMessage] 新消息 ", string(data))
}

type LocationMessage struct {
	MsgId   string                         `json:"msg_id"`
	Event   string                         `json:"event"`
	Content message.LocationMessageRequest `json:"content"`
}

// onLocationMessage 位置消息
func (h *Handler) onLocationMessage(_ context.Context, _ socket.IClient, data []byte) {

	var m LocationMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat onLocationMessage Err: ", err)
		return
	}

	fmt.Println("[onLocationMessage] 新消息 ", string(data))
}

type VoteMessage struct {
	MsgId   string                     `json:"msg_id"`
	Event   string                     `json:"event"`
	Content message.VoteMessageRequest `json:"content"`
}

// onVoteMessage 投票消息
func (h *Handler) onVoteMessage(_ context.Context, _ socket.IClient, data []byte) {

	var m VoteMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat onVoteMessage Err: ", err)
		return
	}

	fmt.Println("[onVoteMessage] 新消息 ", string(data))
}
