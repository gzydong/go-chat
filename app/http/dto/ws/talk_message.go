package ws

type TalkMessageContent struct {
	ReceiverID int              `json:"receiver_id"`
	SenderID   int              `json:"sender_id"`
	TalkType   int              `json:"talk_type"`
	Data       *TalkMessageData `json:"data"`
}

type TalkMessageData struct {
	Avatar      string         `json:"avatar"`
	CodeBlock   *CodeBlockData `json:"code_block"`
	Content     string         `json:"content"`
	CreatedAt   string         `json:"created_at"`
	File        *FileData      `json:"file"`
	Forward     []interface{}  `json:"forward"`
	GroupAvatar string         `json:"group_avatar"`
	GroupName   string         `json:"group_name"`
	ID          int            `json:"id"`
	Invite      *InviteData    `json:"invite"`
	IsRevoke    int            `json:"is_revoke"`
	MsgType     int            `json:"msg_type"`
	Nickname    string         `json:"nickname"`
	ReceiverID  int            `json:"receiver_id"`
	TalkType    int            `json:"talk_type"`
	UserID      int            `json:"user_id"`
}

type CodeBlockData struct {
}

type FileData struct {
}

type InviteData struct {
}
