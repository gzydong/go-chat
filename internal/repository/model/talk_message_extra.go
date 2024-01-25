package model

type Quote struct {
	QuoteId  string `json:"quote_id"`
	Nickname string `json:"nickname"`
	MsgType  int    `json:"msg_type"`
	Content  string `json:"content"`
}

type TalkRecordExtraGroupMember struct {
	UserId   int    `json:"user_id"`  // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
}

// TalkRecordExtraText 文本消息
type TalkRecordExtraText struct {
	Content  string `json:"content"`            // 文本消息
	Mentions []int  `json:"mentions,omitempty"` // @用户ID列表
}

// TalkRecordExtraCode 代码消息
type TalkRecordExtraCode struct {
	Lang string `json:"lang"` // 语言类型
	Code string `json:"code"` // 代码内容
}

// TalkRecordExtraLocation 位置消息
type TalkRecordExtraLocation struct {
	Longitude   string `json:"longitude"`   // 经度
	Latitude    string `json:"latitude"`    // 纬度
	Description string `json:"description"` // 位置描述
}

// TalkRecordExtraForward 转发消息
type TalkRecordExtraForward struct {
	TalkType   int                            `json:"talk_type"`   // 对话类型
	UserId     int                            `json:"user_id"`     // 发送者ID
	ReceiverId int                            `json:"receiver_id"` // 接收者ID（群ID或者好友ID）
	MsgIds     []string                       `json:"msg_ids"`     // 消息列表
	Records    []TalkRecordExtraForwardRecord `json:"records"`     // 消息快照
}

type TalkRecordExtraForwardRecord struct {
	Nickname string `json:"nickname"` // 用户昵称
	Content  string `json:"content"`  // 消息摘要
}

// TalkRecordExtraLogin 登录消息
type TalkRecordExtraLogin struct {
	IP       string `json:"ip"`       // 登录IP
	Address  string `json:"address"`  // 登录地址
	Agent    string `json:"agent"`    // 登录设备
	Platform string `json:"platform"` // 登录平台
	Reason   string `json:"reason"`   // 登录原因
	Datetime string `json:"datetime"` // 登录时间
}

type TalkRecordExtraFile struct {
	Name string `json:"name"` // 文件名称
	Size int    `json:"size"` // 文件大小
	Path string `json:"path"` // 文件路径
}

type TalkRecordExtraImage struct {
	Name   string `json:"name"`   // 图片名称
	Size   int    `json:"size"`   // 图片大小
	Url    string `json:"url"`    // 图片地址
	Width  int    `json:"width"`  // 图片宽度
	Height int    `json:"height"` // 图片高度
}

type TalkRecordExtraAudio struct {
	Name     string `json:"name"`     // 语音名称
	Size     int    `json:"size"`     // 语音大小
	Url      string `json:"url"`      // 语音地址
	Duration int    `json:"duration"` // 语音时长
}

type TalkRecordExtraVideo struct {
	Name     string `json:"name"`     // 视频名称
	Cover    string `json:"cover"`    // 视频封面
	Size     int    `json:"size"`     // 视频大小
	Url      string `json:"url"`      // 视频地址
	Duration int    `json:"duration"` // 视频时长
}

// TalkRecordExtraGroupCreate 创建群消息
type TalkRecordExtraGroupCreate struct {
	OwnerId   int                          `json:"owner_id"`   // 操作人ID
	OwnerName string                       `json:"owner_name"` // 操作人昵称
	Members   []TalkRecordExtraGroupMember `json:"members"`    // 成员列表
}

// TalkRecordExtraGroupJoin 群主邀请加入群消息
type TalkRecordExtraGroupJoin struct {
	OwnerId   int                          `json:"owner_id"`   // 操作人ID
	OwnerName string                       `json:"owner_name"` // 操作人昵称
	Members   []TalkRecordExtraGroupMember `json:"members"`    // 成员列表
}

// TalkRecordExtraTransferGroup 群主转让群消息
type TalkRecordExtraTransferGroup struct {
	OldOwnerId   int    `json:"old_owner_id"`   // 老群主ID
	OldOwnerName string `json:"old_owner_name"` // 老群主昵称
	NewOwnerId   int    `json:"new_owner_id"`   // 新群主ID
	NewOwnerName string `json:"new_owner_name"` // 新群主昵称
}

// TalkRecordExtraGroupMuted 管理员设置群禁言消息
type TalkRecordExtraGroupMuted struct {
	OwnerId   int    `json:"owner_id"`   // 操作人ID
	OwnerName string `json:"owner_name"` // 操作人昵称
}

// TalkRecordExtraGroupCancelMuted 管理员解除群禁言消息
type TalkRecordExtraGroupCancelMuted struct {
	OwnerId   int    `json:"owner_id"`   // 操作人ID
	OwnerName string `json:"owner_name"` // 操作人昵称
}

// TalkRecordExtraGroupMemberMuted 管理员设置群成员禁言消息
type TalkRecordExtraGroupMemberMuted struct {
	OwnerId   int                          `json:"owner_id"`   // 操作人ID
	OwnerName string                       `json:"owner_name"` // 操作人昵称
	Members   []TalkRecordExtraGroupMember `json:"members"`    // 成员列表
}

// TalkRecordExtraGroupMemberCancelMuted 管理员解除群成员禁言消息
type TalkRecordExtraGroupMemberCancelMuted struct {
	OwnerId   int                          `json:"owner_id"`   // 操作人ID
	OwnerName string                       `json:"owner_name"` // 操作人昵称
	Members   []TalkRecordExtraGroupMember `json:"members"`    // 成员列表
}

// TalkRecordExtraGroupDismissed 群主解散群消息
type TalkRecordExtraGroupDismissed struct {
	OwnerId   int    `json:"owner_id"`   // 操作人ID
	OwnerName string `json:"owner_name"` // 操作人昵称
}

// TalkRecordExtraGroupMemberQuit 群成员退出群消息
type TalkRecordExtraGroupMemberQuit struct {
	OwnerId   int    `json:"owner_id"`   // 操作人ID
	OwnerName string `json:"owner_name"` // 操作人昵称
}

// TalkRecordExtraGroupMemberKicked 踢出群成员消息
type TalkRecordExtraGroupMemberKicked struct {
	OwnerId   int                          `json:"owner_id"`   // 操作人ID
	OwnerName string                       `json:"owner_name"` // 操作人昵称
	Members   []TalkRecordExtraGroupMember `json:"members"`    // 成员列表
}

// TalkRecordExtraGroupMessageRevoke 管理员撤回成员消息
type TalkRecordExtraGroupMessageRevoke struct {
	OwnerId     int    `json:"owner_id"`      // 操作人ID
	OwnerName   string `json:"owner_name"`    // 操作人昵称
	RevokeMsgId string `json:"revoke_msg_id"` // 被撤回消息ID
}

// TalkRecordExtraGroupNotice 发布群公告
type TalkRecordExtraGroupNotice struct {
	OwnerId   int    `json:"owner_id"`   // 操作人ID
	OwnerName string `json:"owner_name"` // 操作人昵称
	Title     string `json:"title"`      // 标题
	Content   string `json:"content"`    // 内容
}

type TalkRecordExtraMixedItem struct {
	Type    int    `json:"type"`           // 消息类型, 跟msgtype字段一致
	Content string `json:"content"`        // 消息内容。可包含图片、文字、表情等多种消息。
	Link    string `json:"link,omitempty"` // 图片跳转地址
}

// TalkRecordExtraMixed 图文混合消息
type TalkRecordExtraMixed struct {
	// 消息内容。可包含图片、文字、等消息。
	Items []*TalkRecordExtraMixedItem `json:"items"` // 消息内容。可包含图片、文字、表情等多种消息。
}

type TalkRecordExtraVote struct {
	VoteId int `json:"vote_id"` // 群投票ID
}

// TalkRecordExtraUserShare 用户名片消息
type TalkRecordExtraUserShare struct {
	UserId   int    `json:"user_id"`  // 名片用户ID
	Nickname string `json:"nickname"` // 名片用户昵称
	Avatar   string `json:"avatar"`   // 名片用户头像
	Describe string `json:"describe"` // 群描述
}

// TalkRecordExtraCard 卡片消息
type TalkRecordExtraCard struct {
	Title    string `json:"title"`    // 标题
	Describe string `json:"describe"` // 描述
	Content  string `json:"content"`  // 内容
	Link     string `json:"link"`     // 跳转链接
	Avatar   string `json:"avatar"`   // 头像
	BtnText  string `json:"btn_text"` // 按钮文字
}

// TalkRecordExtraGroupShare 群名片消息
type TalkRecordExtraGroupShare struct {
	GroupId   int    `json:"group_id"`   // 群ID
	GroupName string `json:"group_name"` // 群名
	Describe  string `json:"describe"`   // 群描述
	Avatar    string `json:"avatar"`     // 群头像
	Total     int    `json:"total"`      // 群成员数
}
