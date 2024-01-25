package entity

const (
	SubEventImMessage         = "sub.im.message"          // 对话消息通知
	SubEventImMessageKeyboard = "sub.im.message.keyboard" // 键盘输入事件通知
	SubEventImMessageRevoke   = "sub.im.message.revoke"   // 聊天消息撤销通知
	SubEventContactStatus     = "sub.im.contact.status"   // 用户在线状态通知
	SubEventContactApply      = "sub.im.contact.apply"    // 好友申请消息通知
	SubEventGroupJoin         = "sub.im.group.join"       // 邀请加入群聊通知
	SubEventGroupApply        = "sub.im.group.apply"      // 入群申请通知
)

type SubscribeMessage struct {
	Event   string `json:"event"`   // 事件
	Payload string `json:"payload"` // json 字符串
}

type SubEventImMessagePayload struct {
	TalkMode int    `json:"talk_mode"` // 1 单聊 2 群聊
	Message  string `json:"message"`   // json 字符串
}

type SubEventGroupJoinPayload struct {
	Type    int   `json:"type"` // 1 加入 2 退出
	GroupId int   `json:"group_id"`
	Uids    []int `json:"uids"`
}

type SubEventGroupApplyPayload struct {
	GroupId int `json:"group_id"`
	UserId  int `json:"user_id"`
	ApplyId int `json:"apply_id"`
}

type SubEventContactApplyPayload struct {
	ApplyId int `json:"apply_id"`
	Type    int `json:"type"`
}

type SubEventImMessageKeyboardPayload struct {
	FromId   int `json:"from_id"`
	ToFromId int `json:"to_from_id"`
}

type SubEventContactStatusPayload struct {
	Status int `json:"status"` // 1:上线 2:下线
	UserId int `json:"user_id"`
}
