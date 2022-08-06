package web

type TalkListItem struct {
	Id         int32  `json:"id"`
	TalkType   int32  `json:"talk_type"`
	ReceiverId int32  `json:"receiver_id"`
	IsTop      int32  `json:"is_top"`
	IsDisturb  int32  `json:"is_disturb"`
	IsOnline   int32  `json:"is_online"`
	IsRobot    int32  `json:"is_robot"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	RemarkName string `json:"remark_name"`
	UnreadNum  int32  `json:"unread_num"`
	MsgText    string `json:"msg_text"`
	UpdatedAt  string `json:"updated_at"`
}

// 创建会话列表接口
type (
	CreateTalkRequest struct {
		TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
		ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,min=1" label:"receiver_id"`
	}

	CreateTalkResponse struct {
		Id         int32  `json:"id"`
		TalkType   int32  `json:"talk_type"`
		ReceiverId int32  `json:"receiver_id"`
		IsTop      int32  `json:"is_top"`
		IsDisturb  int32  `json:"is_disturb"`
		IsOnline   int32  `json:"is_online"`
		IsRobot    int32  `json:"is_robot"`
		Name       string `json:"name"`
		Avatar     string `json:"avatar"`
		RemarkName string `json:"remark_name"`
		UnreadNum  int32  `json:"unread_num"`
		MsgText    string `json:"msg_text"`
		UpdatedAt  string `json:"updated_at"`
	}
)

// 获取会话列表接口
type (
	GetTalkListResponse struct {
		Items []*TalkListItem `json:"items"`
	}

	GetTalkListRequest struct{}
)

// 删除会话列表接口
type (
	DeleteTalkListRequest struct {
		Id int `form:"list_id" json:"list_id" binding:"required,numeric" label:"list_id"`
	}

	DeleteTalkListResponse struct{}
)

// 会话置顶接口
type (
	TopTalkListRequest struct {
		Id   int `form:"list_id" json:"list_id" binding:"required,numeric" label:"list_id"`
		Type int `form:"type" json:"type" binding:"required,oneof=1 2" label:"type"`
	}

	TopTalkListResponse struct{}
)

// 会话免打扰接口
type (
	DisturbTalkListRequest struct {
		TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
		ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
		IsDisturb  int `form:"is_disturb" json:"is_disturb" binding:"oneof=0 1" label:"is_disturb"`
	}

	DisturbTalkListResponse struct{}
)

// 获取会话详情接口
type (
	GetTalkRecordsRequest struct {
		TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2"`         // 对话类型
		MsgType    int `form:"msg_type" json:"msg_type" binding:"numeric"`                      // 消息类型
		ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,min=1"` // 接收者ID
		RecordId   int `form:"record_id" json:"record_id" binding:"min=0,numeric"`              // 上次查询的最小消息ID
		Limit      int `form:"limit" json:"limit" binding:"required,numeric,max=100"`           // 数据行数
	}

	GetTalkRecordsResponse struct{}
)

// 清除对话未读数接口
type (
	ClearTalkUnreadNumRequest struct {
		TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
		ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	}

	ClearTalkUnreadNumResponse struct{}
)

// 获取会话转发记录接口
type (
	GetForwardTalkRecordRequest struct {
		RecordId int `form:"record_id" json:"record_id" binding:"min=0,numeric"` // 上次查询的最小消息ID
	}

	GetForwardTalkRecordResponse struct{}
)
