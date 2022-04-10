package entity

// 聊天模式
const (
	ChatPrivateMode = 1 // 私信模式
	ChatGroupMode   = 2 // 群聊模式
)

// WebSocket 消息事件枚举
const (
	EventTalk          = "event_talk"            // 对话消息通知
	EventKeyboard      = "event_keyboard"        // 键盘输入事件通知
	EventOnlineStatus  = "event_online_status"   // 用户在线状态通知
	EventRevokeTalk    = "event_revoke_talk"     // 聊天消息撤销通知
	EventFriendApply   = "event_friend_apply"    // 好友申请消息通知
	EventJoinGroupRoom = "event_join_group_room" // 入群通知
)

// 聊天消息类型
const (
	MsgTypeSysText     = 0  // 系统文本消息
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
