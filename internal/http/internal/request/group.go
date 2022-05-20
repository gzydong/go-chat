package request

type GroupCommonRequest struct {
	GroupId int `form:"group_id" json:"group_id" binding:"required"`
}

type GroupCreateRequest struct {
	Name       string `form:"name" json:"name" binding:"required" label:"name"`
	MembersIds string `form:"ids" json:"ids" binding:"required,ids" label:"ids"`
	Avatar     string `form:"avatar" json:"avatar" label:"avatar"`
	Profile    string `form:"profile" json:"profile" binding:"max=255" label:"profile"`
}

type GroupDismissRequest struct {
	GroupCommonRequest
}

type GroupInviteRequest struct {
	GroupCommonRequest
	Ids string `form:"ids" json:"ids" binding:"required,ids"`
}

type GroupSecedeRequest struct {
	GroupCommonRequest
}

type GroupSettingRequest struct {
	GroupCommonRequest
	GroupName string `form:"group_name" json:"group_name" binding:"required"`
	Avatar    string `form:"avatar" json:"avatar" binding:""`
	Profile   string `form:"profile" json:"profile" binding:"max=255"`
}

type GroupRemoveMembersRequest struct {
	GroupCommonRequest
	MembersIds string `form:"members_ids" json:"members_ids" binding:"required,ids"`
}

type GroupDetailRequest struct {
	GroupCommonRequest
}

type GroupEditRemarkRequest struct {
	GroupCommonRequest
	VisitCard string `form:"visit_card" json:"visit_card" binding:"max=20" label:"visit_card"`
}

type GroupEditNoticeRequest struct {
	GroupCommonRequest
	NoticeId  int    `form:"notice_id" json:"notice_id" binding:"required"`
	Title     string `form:"title" json:"title" binding:"required,max=50"`
	Content   string `form:"content" json:"content" binding:""`
	IsTop     int    `form:"is_top" json:"is_top" binding:"required"`
	IsConfirm string `form:"is_confirm" json:"is_confirm" binding:"required"`
}

type GroupDeleteNoticeRequest struct {
	GroupCommonRequest
	NoticeId int `form:"notice_id" json:"notice_id" binding:"required"`
}

type GetInviteFriendsRequest struct {
	GroupId int `form:"group_id" json:"group_id" binding:"min=0"`
}

type GroupOvertListRequest struct {
	Page int    `form:"page" json:"page" binding:"required"`
	Name string `form:"name" json:"name" binding:"max=50"`
}

type GroupHandoverRequest struct {
	GroupId int `form:"group_id" json:"group_id" binding:"min=1"`
	UserId  int `form:"user_id" json:"user_id" binding:"min=1"`
}

type GroupAssignAdminRequest struct {
	Mode    int `form:"mode" json:"mode" binding:"required,oneof=1 2"`
	GroupId int `form:"group_id" json:"group_id" binding:"min=1"`
	UserId  int `form:"user_id" json:"user_id" binding:"min=1"`
}

type GroupNoSpeakRequest struct {
	Mode    int `form:"mode" json:"mode" binding:"required,oneof=1 2"`
	GroupId int `form:"group_id" json:"group_id" binding:"min=1"`
	UserId  int `form:"user_id" json:"user_id" binding:"min=1"`
}
