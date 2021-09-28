package model

type TalkRecordsCode struct {
	ID        int    `json:"id" grom:"comment:代码块ID"`
	RecordId  int    `json:"record_id" grom:"comment:聊天记录ID"`
	UserId    int    `json:"user_id" grom:"comment:用户ID"`
	CodeLang  string `json:"code_lang" grom:"comment:代码语言"`
	Code      string `json:"code" grom:"comment:代码详情"`
	CreatedAt string `json:"created_at" grom:"comment:创建时间"`
}
