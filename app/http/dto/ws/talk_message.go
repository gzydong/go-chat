package ws

type TalkMessageContent struct {
	ReceiverID int              `json:"receiver_id"`
	SenderID   int              `json:"sender_id"`
	TalkType   int              `json:"talk_type"`
	Data       *TalkMessageItem `json:"data"`
}

type TalkMessageItem struct {
	ID          int           `json:"id"`
	TalkType    int           `json:"talk_type"`
	MsgType     int           `json:"msg_type"`
	UserID      int           `json:"user_id"`
	ReceiverID  int           `json:"receiver_id"`
	Nickname    string        `json:"nickname"`
	Avatar      string        `json:"avatar"`
	GroupName   string        `json:"group_name"`
	GroupAvatar string        `json:"group_avatar"`
	IsRevoke    int           `json:"is_revoke"`
	IsMark      int           `json:"is_mark"`
	IsRead      int           `json:"is_read"`
	Content     string        `json:"content"`
	File        []interface{} `json:"file"`
	CodeBlock   []interface{} `json:"code_block"`
	Forward     []interface{} `json:"forward"`
	Invite      []interface{} `json:"invite"`
	Vote        []interface{} `json:"vote"`
	Login       []interface{} `json:"login"`
	CreatedAt   string        `json:"created_at"`
}

type CodeBlockData struct {
}

type FileData struct {
}

type InviteData struct {
}
