package model

type TalkList struct {
	ID         int    `json:"id"`
	TalkType   string `json:"talk_type"`
	UserId     string `json:"user_id"`
	ReceiverId string `json:"receiver_id"`
	IsDelete   string `json:"is_delete"`
	IsTop      string `json:"is_top"`
	IsRobot    string `json:"is_robot"`
	IsDisturb  string `json:"is_disturb"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
