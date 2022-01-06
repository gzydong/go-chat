package model

type TalkRecordsInvite struct {
	Id            int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                   // 入群或退群通知ID
	RecordId      int    `gorm:"column:record_id;default:0;NOT NULL" json:"record_id"`             // 消息记录ID
	Type          int    `gorm:"column:type;default:1;NOT NULL" json:"type"`                       // 通知类型 （1:入群通知 2:自动退群 3:管理员踢群）
	OperateUserId int    `gorm:"column:operate_user_id;default:0;NOT NULL" json:"operate_user_id"` // 操作人的用户ID（邀请人OR管理员ID）
	UserIds       string `gorm:"column:user_ids;NOT NULL" json:"user_ids"`                         // 用户ID，多个用 , 分割
}
