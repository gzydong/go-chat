package entity

// 聊天模式
const (
	ChatPrivateMode = 1 // 私信模式
	ChatGroupMode   = 2 // 群聊模式
	ChatRoomMode    = 3 // 房间模式
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
	MsgTypeEmoticon    = 11 // 表情消息
)

const (
	BusinessCodeTalk = 101

	BusinessCodeExample = 102
)

// IM消息类型
// 1-999    自定义消息类型
// 1000-1999 系统消息类型
const (
	ChatMsgTypeText     = 1  // 文本消息
	ChatMsgTypeCode     = 2  // 代码消息
	ChatMsgTypeEmoticon = 3  // 表情消息
	ChatMsgTypeImage    = 3  // 图片文件消息
	ChatMsgTypeVoice    = 4  // 语音文件消息
	ChatMsgTypeVideo    = 5  // 视频文件消息
	ChatMsgTypeFile     = 6  // 其它文件消息
	ChatMsgTypeLocation = 7  // 位置消息
	ChatMsgTypeCard     = 8  // 名片消息
	ChatMsgTypeForward  = 9  // 转发消息
	ChatMsgTypeLogin    = 10 // 登录消息

	ChatMsgSysText                   = 1000 // 系统文本消息
	ChatMsgSysGroupCreate            = 1101 // 创建群聊消息
	ChatMsgSysGroupMemberJoin        = 1102 // 加入群聊消息
	ChatMsgSysGroupMemberQuit        = 1103 // 群成员退出群消息
	ChatMsgSysGroupMemberKicked      = 1104 // 踢出群成员消息
	ChatMsgSysGroupMessageRevoke     = 1105 // 管理员撤回成员消息
	ChatMsgSysGroupDismissed         = 1106 // 群解散
	ChatMsgSysGroupMuted             = 1107 // 群禁言
	ChatMsgSysGroupCancelMuted       = 1108 // 群解除禁言
	ChatMsgSysGroupMemberMuted       = 1109 // 群成员禁言
	ChatMsgSysGroupMemberCancelMuted = 1110 // 群成员解除禁言
)

const (
	EventAck = 1000 // ACK消息确认

	EventChatTalkMessage    = 101001 // IM对话消息事件
	EventChatTalkKeyboard   = 101002 // IM键盘输入消息事件
	EventChatTalkRevoke     = 101002 // IM消息撤回事件
	EventChatOnlineStatus   = 101003 // IM在线状态事件
	EventChatContactApply   = 101004 // IM好友申请事件
	EventChatGroupJoinApply = 101005 // IM群加入申请事件
)
