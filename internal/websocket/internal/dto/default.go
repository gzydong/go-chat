package dto

type RevokeTalkMessageContent struct {
	ReceiverID int `json:"receiver_id"`
	RecordID   int `json:"record_id"`
	SenderID   int `json:"sender_id"`
	TalkType   int `json:"talk_type"`
}

type LoginMessageContent struct {
	Status int `json:"status"`
	UserID int `json:"user_id"`
}

type KeyboardMessageContent struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
}

type AckReplyContent struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	MsgID      int `json:"msg_id"`
}

type KeyboardMessage struct {
	Event string                 `json:"event"`
	Data  KeyboardMessageContent `json:"data"`
}

type TalkReadMessage struct {
	Event string `json:"event"`
	Data  struct {
		MsgIds     []int `json:"msg_id"`
		ReceiverId int   `json:"receiver_id"`
	} `json:"data"`
}
