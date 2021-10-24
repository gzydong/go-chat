package request

type GroupNoticeCommonRequest struct {
	GroupId  int `form:"group_id" binding:"required" label:"group_id"`
	NoticeId int `form:"notice_id" binding:"required" label:"notice_id"`
}

type GroupNoticeEditRequest struct {
	GroupId   int    `form:"group_id" binding:"required" label:"group_id"`
	NoticeId  int    `form:"notice_id" label:"notice_id"`
	Title     string `form:"title" binding:"required,max=50" label:"title"`
	Content   string `form:"content" binding:"required,max=65535" label:"content"`
	IsTop     int    `form:"is_top" binding:"oneof=0 1" label:"is_top"`
	IsConfirm int    `form:"is_confirm" binding:"oneof=0 1" label:"is_confirm"`
}

type GroupNoticeListRequest struct {
	GroupId int `form:"group_id" binding:"required" label:"group_id"`
}
