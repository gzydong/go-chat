package dto

type TalkListItem struct {
	Id         int    `json:"id"`
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

type TalkRecordsItem struct {
	Id         int         `json:"id"`
	TalkType   int         `json:"talk_type"`
	MsgType    int         `json:"msg_type"`
	UserId     int         `json:"user_id"`
	ReceiverId int         `json:"receiver_id"`
	Nickname   string      `json:"nickname"`
	Avatar     string      `json:"avatar"`
	IsRevoke   int         `json:"is_revoke"`
	IsMark     int         `json:"is_mark"`
	IsRead     int         `json:"is_read"`
	Content    string      `json:"content,omitempty"`
	File       interface{} `json:"file,omitempty"`
	CodeBlock  interface{} `json:"code_block,omitempty"`
	Forward    interface{} `json:"forward,omitempty"`
	Invite     interface{} `json:"invite,omitempty"`
	Vote       interface{} `json:"vote,omitempty"`
	Login      interface{} `json:"login,omitempty"`
	Location   interface{} `json:"location,omitempty"`
	CreatedAt  string      `json:"created_at"`
}
