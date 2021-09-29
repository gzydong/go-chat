package model

type TalkRecordsLogin struct {
	ID        int    `json:"id" grom:"comment:登录ID"`
	RecordId  int    `json:"record_id" grom:"comment:消息记录ID"`
	UserId    int    `json:"user_id" grom:"comment:用户ID"`
	Ip        string `json:"ip" grom:"comment:登录IP"`
	Platform  string `json:"platform" grom:"comment:登录平台"`
	Agent     string `json:"agent" grom:"comment:设备信息"`
	Address   string `json:"address" grom:"comment:登录地址"`
	Reason    string `json:"reason" grom:"comment:异常信息"`
	CreatedAt string `json:"created_at" grom:"comment:登录时间"`
}
