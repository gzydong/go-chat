package request

type TalkListCreateRequest struct {
	TalkType   int `form:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" binding:"required,numeric" label:"receiver_id"`
}

type TalkListDeleteRequest struct {
	Id int `form:"list_id" binding:"required,numeric" label:"list_id"`
}

type TalkListTopRequest struct {
	Id   int `form:"list_id" binding:"required,numeric" label:"list_id"`
	Type int `form:"type" binding:"required,oneof=1 2" label:"type"`
}

type TalkListDisturbRequest struct {
	TalkType   int `form:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	IsDisturb  int `form:"is_disturb" binding:"oneof=0 1" label:"is_disturb"`
}
