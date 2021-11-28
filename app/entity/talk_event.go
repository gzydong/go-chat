package entity

// WebSocket 消息事件枚举

const (
	// 对话消息通知 - 事件名
	EventTalk = "event_talk"

	// 键盘输入事件通知 - 事件名
	EventKeyboard = "event_keyboard"

	// 用户在线状态通知 - 事件名
	EventOnlineStatus = "event_online_status"

	// 聊天消息撤销通知 - 事件名
	EventRevokeTalk = "event_revoke_talk"

	// 好友申请消息通知 - 事件名
	EventFriendApply = "event_friend_apply"

	EventJoinGroupRoom = "event_join_group_room"
)
