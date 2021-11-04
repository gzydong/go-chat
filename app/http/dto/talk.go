package dto

type TalkListItem struct {
	ID         int    `json:"id"`
	TalkType   int    `json:"talk_type"`
	ReceiverId int    `json:"receiver_id"`
	IsTop      int    `json:"is_top"`
	IsDisturb  int    `json:"is_disturb"`
	IsOnline   int    `json:"is_online"`
	IsRobot    int    `json:"is_robot"`
	Avatar     string `json:"avatar"`
	Name       string `json:"name"`
	RemarkName string `json:"remark_name"`
	UnreadNum  int    `json:"unread_num"`
	MsgText    string `json:"msg_text"`
	UpdatedAt  string `json:"updated_at"`
}

type TalkMessageContent struct {
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
