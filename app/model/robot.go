package model

type Robot struct {
	ID        int    `json:"id" grom:"comment:机器人ID"`
	UserId    int    `json:"user_id" grom:"comment:关联用户ID"`
	RobotName string `json:"robot_name" grom:"comment:机器人名称"`
	Describe  string `json:"describe" grom:"comment:描述信息"`
	Logo      string `json:"logo" grom:"comment:机器人logo"`
	IsTalk    int    `json:"is_talk" grom:"comment:可发送消息"`
	Status    int    `json:"status" grom:"comment:状态"`
	Type      int    `json:"type" grom:"comment:机器人类型"`
	CreatedAt string `json:"created_at" grom:"comment:创建时间"`
	UpdatedAt string `json:"updated_at" grom:"comment:更新时间"`
}
