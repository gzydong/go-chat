package model

import "time"

type TalkRecordsCode struct {
	ID        int       `json:"-" grom:"comment:代码块ID"`
	RecordId  int       `json:"-" grom:"comment:聊天记录ID"`
	UserId    int       `json:"-" grom:"comment:用户ID"`
	CodeLang  string    `json:"code_lang" grom:"comment:代码语言"`
	Code      string    `json:"code" grom:"comment:代码详情"`
	CreatedAt time.Time `json:"-" grom:"comment:创建时间"`
}
