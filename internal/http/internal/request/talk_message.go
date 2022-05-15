package request

type TextMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	Text       string `form:"text" json:"text" binding:"required,max=3000" label:"text"`
}

type CodeMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	Lang       string `form:"lang" json:"lang" binding:"required"`
	Code       string `form:"code" json:"code" binding:"required,max=3000"`
}

type ImageMessageRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
}

type FileMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	UploadId   string `form:"upload_id" json:"upload_id" binding:"required"`
}

type VoteMessageRequest struct {
	ReceiverId int      `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	Mode       int      `form:"mode" json:"mode" binding:"oneof=0 1"`
	Anonymous  int      `form:"anonymous" json:"anonymous" binding:"oneof=0 1"`
	Title      string   `form:"title" json:"title" binding:"required"`
	Options    []string `form:"options" json:"options"`
}

type VoteMessageHandleRequest struct {
	RecordId int    `form:"record_id" json:"record_id" binding:"required,gt=0"`
	Options  string `form:"options" json:"options" binding:"required"`
}

type EmoticonMessageRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	EmoticonId int `form:"emoticon_id" json:"emoticon_id" binding:"required,numeric,gt=0"`
}

type ForwardMessageRequest struct {
	TalkType        int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId      int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	ForwardMode     int    `form:"forward_mode" json:"forward_mode" binding:"required,oneof=1 2"`
	RecordsIds      string `form:"records_ids" json:"records_ids" binding:"required,ids"`
	ReceiveUserIds  string `form:"receive_user_ids" json:"receive_user_ids" binding:"ids"`
	ReceiveGroupIds string `form:"receive_group_ids" json:"receive_group_ids" binding:"ids"`
}

type CardMessageRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
}

type LocationMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	Longitude  string `form:"longitude" json:"longitude" binding:"required,numeric"`
	Latitude   string `form:"latitude" json:"latitude" binding:"required,numeric"`
}

type RevokeMessageRequest struct {
	RecordId int `form:"record_id" json:"record_id" binding:"required,numeric,gt=0" label:"record_id"`
}

type CollectMessageRequest struct {
	RecordId int `form:"record_id" json:"record_id" binding:"required,numeric,gt=0" label:"record_id"`
}

type DeleteMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	RecordIds  string `form:"record_id" json:"record_id" binding:"required,ids" label:"record_id"`
}
