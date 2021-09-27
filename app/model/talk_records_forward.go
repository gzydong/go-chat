package model

type TalkRecordsForward struct {
	ID        int    `json:"id" grom:"comment:转发ID"`
	RecordId  int    `json:"record_id" grom:"comment:聊天记录ID"`
	UserId    int    `json:"user_id" grom:"comment:用户ID"`
	RecordsId int    `json:"records_id" grom:"comment:聊天记录ID，多个用英文','拼接"`
	Text      string `json:"text" grom:"comment:缓存信息"`
	CreatedAt string `json:"created_at" grom:"comment:转发时间"`
}
