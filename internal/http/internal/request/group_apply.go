package request

type GroupApplyCommonRequest struct {
	ApplyId int `form:"apply_id" json:"apply_id" binding:"required" label:"apply_id"`
}

type GroupApplyCreateRequest struct {
	GroupId int    `form:"group_id" json:"group_id" binding:"required" label:"group_id"`
	Remark  string `form:"remark" json:"remark" binding:"required,max=255" label:"remark"`
}

type GroupApplyDeleteRequest struct {
	GroupApplyCommonRequest
}

type GroupApplyAgreeRequest struct {
	GroupApplyCommonRequest
}

type GroupApplyListRequest struct {
	GroupId int `form:"group_id" json:"group_id" binding:"required" label:"group_id"`
}
