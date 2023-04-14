package entity

// 聊天模式
const (
	ChatPrivateMode = 1 // 私信模式
	ChatGroupMode   = 2 // 群聊模式
	ChatRoomMode    = 3 // 房间模式
)

const (
	SubEventImMessage         = "sub.im.message"          // 对话消息通知
	SubEventImMessageKeyboard = "sub.im.message.keyboard" // 键盘输入事件通知
	SubEventImMessageRevoke   = "sub.im.message.revoke"   // 聊天消息撤销通知
	SubEventImMessageRead     = "sub.im.message.read"     // 对话消息读事件
	SubEventContactStatus     = "sub.im.contact.status"   // 用户在线状态通知
	SubEventContactApply      = "sub.im.contact.apply"    // 好友申请消息通知
	SubEventGroupJoin         = "sub.im.group.join"       // 邀请加入群聊通知

	PushEventImMessage         = "im.message"          // 对话消息推送
	PushEventImMessageKeyboard = "im.message.keyboard" // 键盘输入事件推送
	PushEventImMessageRead     = "im.message.read"     // 对话消息读事件推送
	PushEventImMessageRevoke   = "im.message.revoke"   // 聊天消息撤销推送
	PushEventContactApply      = "im.contact.apply"    // 好友申请消息推送
	PushEventContactStatus     = "im.contact.status"   // 用户在线状态推送
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
	ChatMsgTypeImage    = 3  // 图片文件
	ChatMsgTypeVoice    = 4  // 语音文件
	ChatMsgTypeVideo    = 5  // 视频文件
	ChatMsgTypeFile     = 6  // 其它文件
	ChatMsgTypeLocation = 7  // 位置消息
	ChatMsgTypeCard     = 8  // 名片消息
	ChatMsgTypeForward  = 9  // 转发消息
	ChatMsgTypeLogin    = 10 // 登录消息
	ChatMsgTypeVote     = 11 // 投票消息

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
	ChatMsgSysGroupNoticeCreate      = 1111 // 创建群公告
	ChatMsgSysGroupNoticeUpdate      = 1112 // 更新群公告
)

const (
	EventChatTalkMessage    = 101001 // IM对话消息事件
	EventChatTalkKeyboard   = 101002 // IM键盘输入消息事件
	EventChatTalkRevoke     = 101002 // IM消息撤回事件
	EventChatOnlineStatus   = 101003 // IM在线状态事件
	EventChatContactApply   = 101004 // IM好友申请事件
	EventChatGroupJoinApply = 101005 // IM群加入申请事件
)
