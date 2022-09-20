package web

type SendBaseMessageRequest struct {
	Type     int       `json:"type" binding:"required,gt=0"`
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// Receiver 接受者信息
type Receiver struct {
	TalkType   int `json:"talk_type" binding:"required,gt=0"`   // 对话类型 1:私聊 2:群聊
	ReceiverId int `json:"receiver_id" binding:"required,gt=0"` // 好友ID或群ID
}

// SendTextRequest 发送文本消息
type SendTextRequest struct {
	Type    int    `json:"type" binding:"required,gt=0"`
	Content string `json:"content" binding:"required"` // 文本信息
	Mention struct {
		All  int      `json:"all"`  // 是否所有成员
		Uids []string `json:"uids"` // 指定成员
	} `json:"mention"` // @ 信息
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// SendImageRequest 发送图片消息
type SendImageRequest struct {
	Type     int       `json:"type" binding:"required,gt=0"`
	Url      string    `json:"url" binding:"required"`
	Width    int       `json:"width" binding:"required,gt=0"`
	Height   int       `json:"height" binding:"required,gt=0"`
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// SendVoiceRequest 发送语音消息
type SendVoiceRequest struct {
	Type     int       `json:"type" binding:"required,gt=0"`
	Url      string    `json:"url" binding:"required"`
	Duration int       `json:"duration" binding:"required,gt=0"`
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// SendFileRequest 发送文件消息
type SendFileRequest struct {
	Type     int       `json:"type" binding:"required,gt=0"`
	Url      string    `json:"url" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Size     int       `json:"size" binding:"required,gt=0"`
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// SendCodeRequest 发送代码消息
type SendCodeRequest struct {
	Type     int       `json:"type" binding:"required,gt=0"`
	Lang     string    `json:"lang" binding:"required"` // 代码类型
	Code     string    `json:"code" binding:"required"` // 代码信息
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// SendForwardRequest 发送转发消息
type SendForwardRequest struct {
	Type       int   `json:"type" binding:"required,gt=0"`
	Mode       int   `json:"mode" binding:"required,oneof=1 2"` // 转发模式
	MessageIds []int `json:"message_ids" binding:"required"`    // 消息ID
	Receiver   struct {
		GroupIds []int `json:"gids"` // 群ID
		Uids     []int `json:"uids"` // 用户ID
	} `json:"mention"`
}

// SendLocationRequest 发送位置消息
type SendLocationRequest struct {
	Type        int       `json:"type" binding:"required,gt=0"`
	Longitude   string    `json:"longitude" binding:"required,numeric"` // 地理位置 经度
	Latitude    string    `json:"latitude" binding:"required,numeric"`  // 地理位置 纬度
	Description string    `json:"description"`                          // 位置描述
	Receiver    *Receiver `json:"receiver" binding:"required"`
}

// SendEmoticonRequest 发送表情包消息
type SendEmoticonRequest struct {
	Type       int       `json:"type" binding:"required,gt=0"`
	EmoticonId int       `json:"emoticon_id" binding:"required,gt=0"` // 表情包ID
	Receiver   *Receiver `json:"receiver" binding:"required"`
}

// SendVoteMessageRequest 发送投票消息
type SendVoteMessageRequest struct {
	Type      int       `json:"type" binding:"required,gt=0"`
	Title     string    `form:"title" json:"title" binding:"required"`          // 标题
	Mode      int       `form:"mode" json:"mode" binding:"oneof=0 1"`           // 投票模式
	Anonymous int       `form:"anonymous" json:"anonymous" binding:"oneof=0 1"` // 匿名投票
	Options   []string  `form:"options" json:"options"`                         // 投票选项
	Receiver  *Receiver `json:"receiver" binding:"required"`
}
