package model

import "time"

type TalkRecordsDelete struct {
	ID        int       `json:"id" grom:"comment:代码块ID"`
	RecordId  int       `json:"record_id" grom:"comment:聊天记录ID"`
	UserId    int       `json:"user_id" grom:"comment:用户ID"`
	CreatedAt time.Time `json:"created_at" grom:"comment:删除时间"`
}
