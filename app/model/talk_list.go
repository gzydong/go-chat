package model

import "time"

type TalkList struct {
	ID         int       `json:"id" grom:"comment:聊天列表ID"`
	TalkType   int       `json:"talk_type" grom:"comment:聊天类型"`
	UserId     int       `json:"user_id" grom:"comment:用户ID或消息发送者ID"`
	ReceiverId int       `json:"receiver_id" grom:"comment:接收者ID"`
	IsDelete   int       `json:"is_delete" grom:"comment:是否删除"`
	IsTop      int       `json:"is_top" grom:"comment:是否置顶"`
	IsRobot    int       `json:"is_robot" grom:"comment:消息免打扰"`
	IsDisturb  int       `json:"is_disturb" grom:"comment:是否机器人"`
	CreatedAt  time.Time `json:"created_at" grom:"comment:创建时间"`
	UpdatedAt  time.Time `json:"updated_at" grom:"comment:更新时间"`
}

type SearchTalkList struct {
	ID          int       `json:"id" grom:"comment:聊天列表ID"`
	TalkType    int       `json:"talk_type" grom:"comment:聊天类型"`
	ReceiverId  int       `json:"receiver_id" grom:"comment:接收者ID"`
	IsDelete    int       `json:"is_delete" grom:"comment:是否删除"`
	IsTop       int       `json:"is_top" grom:"comment:是否置顶"`
	IsRobot     int       `json:"is_robot" grom:"comment:消息免打扰"`
	IsDisturb   int       `json:"is_disturb" grom:"comment:是否机器人"`
	UpdatedAt   time.Time `json:"updated_at" grom:"comment:更新时间"`
	UserAvatar  string    `json:"user_avatar"`
	Nickname    string    `json:"nickname"`
	GroupName   string    `json:"group_name"`
	GroupAvatar string    `json:"group_avatar"`
}
