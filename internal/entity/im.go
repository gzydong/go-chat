package entity

// IM 渠道分组(用于业务划分，业务间相互隔离)
const (
	// ImChannelChat 默认分组
	ImChannelChat    = "chat"    // im.Sessions.Chat.Name()
	ImChannelExample = "example" // im.Sessions.Example.Name()
)

const (
	// ImTopicChat 默认渠道消息订阅
	ImTopicChat        = "im:message:chat:all"
	ImTopicChatPrivate = "im:message:chat:%s"

	// ImTopicExample Example渠道消息订阅
	ImTopicExample        = "im:message:example:all"
	ImTopicExamplePrivate = "im:message:example:%s"
)

type ImMessagePayload struct {
	TalkMode int `json:"talk_mode"`  // 对话类型[1:私信;2:群聊;]
	FromId   int `json:"from_id"`    // 发送者用户ID
	ToFromId int `json:"to_from_id"` // 接收者ID[好友ID或者群ID]
	Body     any `json:"body"`       // 私信消息或群聊消息
}

type ImMessagePayloadBody struct {
	MsgId     string `json:"msg_id"`
	Sequence  int    `json:"sequence"`
	MsgType   int    `json:"msg_type"`
	UserId    int    `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsRevoked int    `json:"is_revoked"`
	SendTime  string `json:"send_time"`
	Extra     any    `json:"extra"` // 额外参数
	Quote     any    `json:"quote"` // 额外参数
}

// ImContactApplyPayload
// im.contact.apply
type ImContactApplyPayload struct {
	UserId    int    `json:"user_id"`
	Nickname  string `json:"nickname"`
	Remark    string `json:"remark"`
	ApplyTime string `json:"apply_time"`
}

// ImContactApplyResultPayload
// im.contact.apply_result
type ImContactApplyResultPayload struct {
	// 同意添加好友的用户ID
	UserId int `json:"user_id"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 申请备注
	ApplyResult string `json:"apply_result"`
	// 操作时间
	OperateTime string `json:"operate_time"`
}

type ImGroupApplyPayload struct {
	GroupId   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	UserId    int    `json:"user_id"`
	Nickname  string `json:"nickname"`
	Remark    string `json:"remark"`
	ApplyTime string `json:"apply_time"`
}

type ImMessageKeyboardPayload struct {
	FromId   int `json:"from_id"`
	ToFromId int `json:"to_from_id"`
}
