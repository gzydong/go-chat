package request

type SetSystemEmoticonRequest struct {
	EmoticonId int `form:"emoticon_id" binding:"required,numeric" label:"emoticon_id"`
	Type       int `form:"type" binding:"required,oneof=1 2" label:"type"`
}

type DeleteCollectRequest struct {
	Ids string `form:"ids" binding:"required,ids" label:"ids"`
}
