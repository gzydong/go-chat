package request

type ContactEditRemarkRequest struct {
	Remarks  string `form:"remarks" json:"remarks" label:"remarks"`
	FriendId int    `form:"friend_id" json:"friend_id" binding:"required" label:"friend_id"`
}

type ContactDeleteRequest struct {
	FriendId int `form:"friend_id" json:"friend_id" binding:"required" label:"friend_id"`
}

type ContactDetailRequest struct {
	UserId int `form:"user_id" json:"user_id" binding:"required,min=1" label:"user_id"`
}

type ContactSearchRequest struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required" label:"mobile"`
}

type ContactApplyCreateRequest struct {
	Remarks  string `form:"remark" json:"remark" binding:"required" label:"remark"`
	FriendId int    `form:"friend_id" json:"friend_id" binding:"required" label:"friend_id"`
}

type ContactApplyAcceptRequest struct {
	Remarks string `form:"remark" json:"remark" binding:"required" label:"remark"`
	ApplyId int    `form:"apply_id" json:"apply_id" binding:"required" label:"apply_id"`
}

type ContactApplyDeclineRequest struct {
	Remarks string `form:"remark" json:"remark" binding:"required" label:"remark"`
	ApplyId int    `form:"apply_id" json:"apply_id" binding:"required" label:"apply_id"`
}
