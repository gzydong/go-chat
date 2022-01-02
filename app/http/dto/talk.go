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
