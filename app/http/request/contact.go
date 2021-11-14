package request

type ContactEditRemarkRequest struct {
	Remarks  string `form:"remarks" json:"remarks" binding:"required" label:"remarks"`
	FriendId int    `form:"friend_id" json:"friend_id" binding:"required" label:"friend_id"`
}

type ContactDeleteRequest struct {
	FriendId int `form:"friend_id" json:"friend_id" binding:"required" label:"friend_id"`
}

type ContactDetailRequest struct {
	UserId int `form:"user_id" json:"user_id" binding:"required,min=1" label:"user_id"`
}

type ContactSearchRequest struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,len=11" label:"mobile"`
}
