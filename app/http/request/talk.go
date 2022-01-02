package request

type TalkListCreateRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,min=1" label:"receiver_id"`
}

type TalkListDeleteRequest struct {
	Id int `form:"list_id" json:"list_id" binding:"required,numeric" label:"list_id"`
}

type TalkListTopRequest struct {
	Id   int `form:"list_id" json:"list_id" binding:"required,numeric" label:"list_id"`
	Type int `form:"type" json:"type" binding:"required,oneof=1 2" label:"type"`
}

type TalkListDisturbRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	IsDisturb  int `form:"is_disturb" json:"is_disturb" binding:"oneof=0 1" label:"is_disturb"`
}

type TalkRecordsRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2"`         // 对话类型
	MsgType    int `form:"msg_type" json:"msg_type" binding:"numeric"`                      // 消息类型
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,min=1"` // 接收者ID
	RecordId   int `form:"record_id" json:"record_id" binding:"min=0,numeric"`              // 上次查询的最小消息ID
	Limit      int `form:"limit" json:"limit" binding:"required,numeric,max=100"`           // 数据行数
}

type TalkUnReadRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
}

type TalkForwardRecordsRequest struct {
	RecordId int `form:"record_id" json:"record_id" binding:"min=0,numeric"` // 上次查询的最小消息ID
}
