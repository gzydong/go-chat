package request

type ArticleClassEditRequest struct {
	ClassId   int    `form:"class_id" json:"class_id"`
	ClassName string `form:"class_name" json:"class_name" binding:"required" label:"class_name"`
}

type ArticleClassDeleteRequest struct {
	ClassId int `form:"class_id" json:"class_id" binding:"required"`
}

type ArticleTagEditRequest struct {
	TagId   int    `form:"tag_id" json:"tag_id"`
	TagName string `form:"tag_name" json:"tag_name" binding:"required" label:"tag_name"`
}

type ArticleTagDeleteRequest struct {
	TagId int `form:"tag_id" json:"tag_id" binding:"required"`
}
