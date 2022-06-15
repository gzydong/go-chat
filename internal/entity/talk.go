package entity

// 聊天模式
const (
	ChatPrivateMode = 1 // 私信模式
	ChatGroupMode   = 2 // 群聊模式
)

// WebSocket 消息事件枚举
const (
	EventTalk          = "event_talk"            // 对话消息通知
	EventTalkKeyboard  = "event_talk_keyboard"   // 键盘输入事件通知
	EventTalkRevoke    = "event_talk_revoke"     // 聊天消息撤销通知
	EventTalkJoinGroup = "event_talk_join_group" // 邀请加入群聊通知
	EventTalkRead      = "event_talk_read"       // 对话消息读事件
	EventOnlineStatus  = "event_login"           // 用户在线状态通知
	EventContactApply  = "event_contact_apply"   // 好友申请消息通知
)

// 聊天消息类型
const (
	MsgTypeSystemText  = 0  // 系统文本消息
	MsgTypeText        = 1  // 文本消息
	MsgTypeFile        = 2  // 文件消息
	MsgTypeForward     = 3  // 会话消息
	MsgTypeCode        = 4  // 代码消息
	MsgTypeVote        = 5  // 投票消息
	MsgTypeGroupNotice = 6  // 群组公告
	MsgTypeFriendApply = 7  // 好友申请
	MsgTypeLogin       = 8  // 登录通知
	MsgTypeGroupInvite = 9  // 入群退群消息
	MsgTypeLocation    = 10 // 位置消息
)
