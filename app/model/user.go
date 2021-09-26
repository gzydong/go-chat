package model

type User struct {
	ID         int    `json:"id"`
	Mobile     string `json:"mobile"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Gender     int    `json:"gender"`
	Password   string `json:"password"`
	InviteCode string `json:"invite_code"`
	Motto      string `json:"motto"`
	Email      string `json:"email"`
	IsRobot    int    `json:"is_robot"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
