package model

type User struct {
	ID        int    `json:"id" grom:"comment:用户ID"`
	Mobile    string `json:"mobile" grom:"comment:用户昵称"`
	Nickname  string `json:"nickname" grom:"comment:登录手机号"`
	Avatar    string `json:"avatar" grom:"comment:邮箱地址"`
	Gender    int    `json:"gender" grom:"comment:登录密码"`
	Password  string `json:"password" grom:"comment:头像"`
	Motto     string `json:"motto" grom:"comment:性别"`
	Email     string `json:"email" grom:"comment:座右铭"`
	IsRobot   int    `json:"is_robot" grom:"comment:是否机器人"`
	CreatedAt string `json:"created_at" grom:"comment:注册时间"`
	UpdatedAt string `json:"updated_at" grom:"comment:更新时间"`
}
