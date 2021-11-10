package model

import "time"

type TalkRecordsLocation struct {
	ID        int       `json:"-" grom:"comment:自增ID"`
	RecordId  int       `json:"-" grom:"comment:聊天记录ID"`
	UserId    int       `json:"-" grom:"comment:用户ID"`
	Longitude string    `json:"longitude" grom:"comment:经度"`
	Latitude  string    `json:"latitude" grom:"comment:纬度"`
	CreatedAt time.Time `json:"-" grom:"comment:创建时间"`
}
