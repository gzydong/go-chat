package model

type TalkRecordsInvite struct {
	ID            int    `json:"id" grom:"comment:入群或退群通知ID"`
	RecordId      int    `json:"record_id" grom:"comment:消息记录ID"`
	Type          int    `json:"type" grom:"comment:通知类通知类型[1:入群通知;2:自动退群;3:管理员踢群]型"`
	OperateUserId int    `json:"operate_user_id" grom:"comment:操作人的用户ID[邀请人OR管理员ID]"`
	UserIds       string `json:"user_ids" grom:"comment:用户ID(多个用 , 分割)"`
}
