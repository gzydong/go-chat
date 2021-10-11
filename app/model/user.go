package model

type User struct {
	ID        int    `json:"id" grom:"comment:用户ID"`
	Mobile    string `json:"mobile" grom:"comment:手机账号"`
	Nickname  string `json:"nickname" grom:"comment:用户昵称"`
	Avatar    string `json:"avatar" grom:"comment:头像"`
	Gender    int    `json:"gender" grom:"comment:性别"`
	Password  string `json:"-" grom:"comment:账号密码"`
	Motto     string `json:"motto" grom:"comment:座右铭"`
	Email     string `json:"email" grom:"comment:绑定邮箱"`
	IsRobot   int    `json:"is_robot" grom:"comment:是否机器人"`
	CreatedAt string `json:"created_at" grom:"comment:注册时间"`
	UpdatedAt string `json:"updated_at" grom:"comment:更新时间"`
}
