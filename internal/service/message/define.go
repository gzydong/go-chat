package message

type CreatePrivateSysMessageOption struct {
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	Content  string `json:"content"`    // 消息内容
}

type CreatePrivateMessageOption struct {
	MsgType  int    `json:"msg_type"`   // 消息类型，1-文本消息，2-图片消息，3-语音消息，4-视频消息，5-文件消息，6-链接消息，7-小程序消息
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	QuoteId  string `json:"quote_id"`   // 引用消息id
	Extra    string `json:"extra"`      // 扩展字段
}

type CreateGroupMessageOption struct {
	MsgType  int    `json:"msg_type"`   // 消息类型，1-文本消息，2-图片消息，3-语音消息，4-视频消息，5-文件消息，6-链接消息，7-小程序消息
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	QuoteId  string `json:"quote_id"`   // 引用消息id
	Extra    string `json:"extra"`      // 扩展字段
}

type CreateGroupSysMessageOption struct {
	GroupId int    `json:"group_id"`
	Content string `json:"content"` // 消息内容
}

type CreateMessageOption struct {
	TalkMode int    `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	MsgType  int    `json:"msg_type"`   // 消息类型
	QuoteId  string `json:"quote_id"`   // 引用消息id
	Extra    string `json:"extra"`      // 扩展字段
}

type CreateLoginMessageOption struct {
	UserId   int    `json:"user_id"`  // 用户ID
	Ip       string `json:"ip"`       // IP地址
	Address  string `json:"address"`  // 地址
	Platform string `json:"platform"` // 平台
	Agent    string `json:"agent"`    // 浏览器
	Reason   string `json:"reason"`   // 登录原因
	LoginAt  string `json:"login_at"` // 登录时间
}

type CreateTextMessage struct {
	TalkMode int    `json:"talk_mode"`          // 发送模式，1-单聊，2-群聊
	FromId   int    `json:"from_id"`            // 发送者
	ToFromId int    `json:"to_from_id"`         // 接受者(好友ID或者群组ID)
	Content  string `json:"content"`            // 消息内容
	QuoteId  string `json:"quote_id"`           // 引用消息id
	Mentions []int  `json:"mentions,omitempty"` // @用户ID列表
}

type CreateImageMessage struct {
	TalkMode int    `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	QuoteId  string `json:"quote_id"`   // 引用消息id
	Url      string `json:"url"`        // 图片地址
	Width    int    `json:"width"`      // 图片宽度
	Height   int    `json:"height"`     // 图片高度
	Size     int    `json:"size"`       // 图片大小
}

type CreateVoiceMessage struct {
	TalkMode int    `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	Url      string `json:"url"`        // 语音地址
	Duration int    `json:"duration"`   // 语音时长
	Size     int    `json:"size"`       // 语音大小
}

type CreateVideoMessage struct {
	TalkMode int    `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	Url      string `json:"url"`        // 视频地址
	Duration int    `json:"duration"`   // 视频时长
	Size     int    `json:"size"`       // 视频大小
	Cover    string `json:"cover"`      // 视频封面
}

type CreateFileMessage struct {
	TalkMode int    `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	UploadId string `json:"upload_id"`  // 文件上传ID
}

type CreateCodeMessage struct {
	TalkMode int    `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int    `json:"from_id"`    // 发送者
	ToFromId int    `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	Code     string `json:"code"`       // 代码内容
	Lang     string `json:"lang"`       // 代码语言
}

type CreateVoteMessage struct {
	TalkMode int `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int `json:"from_id"`    // 发送者
	ToFromId int `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	VoteId   int `json:"vote_id"`    // 投票id
}

type CreateEmoticonMessage struct {
	TalkMode   int `json:"talk_mode"`   // 发送模式，1-单聊，2-群聊
	FromId     int `json:"from_id"`     // 发送者
	ToFromId   int `json:"to_from_id"`  // 接受者(好友ID或者群组ID)
	EmoticonId int `json:"emoticon_id"` // 表情ID
}

type CreateForwardMessage struct {
	TalkMode int      `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int      `json:"from_id"`    // 发送者
	ToFromId int      `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	Action   int      `json:"action"`     // 转发方式 1:合并转发 2:逐条转发
	MsgIds   []string `json:"message_ids"`
	Gids     []int    `json:"gids"` // 群ID列表
	Uids     []int    `json:"uids"` // 好友ID列表
	UserId   int      `json:"user_id"`
}

type CreateLocationMessage struct {
	TalkMode    int    `json:"talk_mode"`   // 发送模式，1-单聊，2-群聊
	FromId      int    `json:"from_id"`     // 发送者
	ToFromId    int    `json:"to_from_id"`  // 接受者(好友ID或者群组ID)
	Longitude   string `json:"longitude"`   // 地理位置 经度
	Latitude    string `json:"latitude"`    // 地理位置 纬度
	Description string `json:"description"` // 位置描述
}

type CreateBusinessCardMessage struct {
	TalkMode int `json:"talk_mode"`  // 发送模式，1-单聊，2-群聊
	FromId   int `json:"from_id"`    // 发送者
	ToFromId int `json:"to_from_id"` // 接受者(好友ID或者群组ID)
	UserId   int `json:"user_id"`    // 名片用户ID
}

type CreateMixedMessage struct {
	TalkMode    int                      `json:"talk_mode"`          // 发送模式，1-单聊，2-群聊
	FromId      int                      `json:"from_id"`            // 发送者
	ToFromId    int                      `json:"to_from_id"`         // 接受者(好友ID或者群组ID)
	QuoteId     string                   `json:"quote_id"`           // 引用消息id
	Mentions    []int                    `json:"mentions,omitempty"` // @用户ID列表
	MessageList []CreateMixedMessageItem `json:"message_list"`       // 消息列表
}

type CreateMixedMessageItem struct {
	Type    int    `json:"type"`    // 消息类型
	Content string `json:"content"` // 消息内容
}
