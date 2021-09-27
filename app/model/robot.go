package model

type Robot struct {
	ID        int    `json:"id"`
	UserId    int    `json:"user_id"`
	RobotName string `json:"robot_name"`
	Describe  string `json:"describe"`
	Logo      string `json:"logo"`
	IsTalk    int    `json:"is_talk"`
	Status    int    `json:"status"`
	Type      int    `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
