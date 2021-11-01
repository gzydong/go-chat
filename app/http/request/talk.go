package request

type TalkListCreateRequest struct {
	TalkType   int `form:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" binding:"required,numeric" label:"receiver_id"`
}
