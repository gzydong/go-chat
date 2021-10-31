package model

import "time"

type TalkRecordsLocation struct {
	ID        int       `json:"id" grom:"comment:自增ID"`
	RecordId  int       `json:"record_id" grom:"comment:聊天记录ID"`
	UserId    int       `json:"user_id" grom:"comment:用户ID"`
	Longitude string    `json:"longitude" grom:"comment:经度"`
	Latitude  string    `json:"latitude" grom:"comment:纬度"`
	CreatedAt time.Time `json:"created_at" grom:"comment:创建时间"`
}
