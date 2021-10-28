package request

type GroupCommonRequest struct {
	GroupId int `form:"group_id" binding:"required"`
}

type GroupCreateRequest struct {
	Name       string `form:"name" binding:"required" label:"name"`
	MembersIds string `form:"ids" binding:"required,ids" label:"ids"`
	Avatar     string `form:"avatar" label:"avatar"`
	Profile    string `form:"profile" binding:"max=255" label:"profile"`
}

type GroupDismissRequest struct {
	GroupCommonRequest
}

type GroupInviteRequest struct {
	GroupCommonRequest
	Ids string `form:"ids" binding:"required,ids"`
}

type GroupSecedeRequest struct {
	GroupCommonRequest
}

type GroupSettingRequest struct {
	GroupCommonRequest
	GroupName string `form:"group_name" binding:"required"`
	Avatar    string `form:"avatar" binding:""`
	Profile   string `form:"profile" binding:"max=255"`
}

type GroupRemoveMembersRequest struct {
	GroupCommonRequest
	MembersIds string `form:"ids" binding:"required,ids"`
}

type GroupDetailRequest struct {
	GroupCommonRequest
}

type GroupEditCardRequest struct {
	GroupCommonRequest
	VisitCard string `form:"visit_card" binding:"required,max=20" label:"visit_card"`
}

type GroupEditNoticeRequest struct {
	GroupCommonRequest
	NoticeId  int    `form:"notice_id" binding:"required"`
	Title     string `form:"title" binding:"required,max=50"`
	Content   string `form:"content" binding:""`
	IsTop     int    `form:"is_top" binding:"required"`
	IsConfirm string `form:"is_confirm" binding:"required"`
}

type GroupDeleteNoticeRequest struct {
	GroupCommonRequest
	NoticeId int `form:"notice_id" binding:"required"`
}
