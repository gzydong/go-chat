package request

type ArticleEditRequest struct {
	ArticleId int    `form:"article_id" json:"article_id"`
	ClassId   int    `form:"class_id" json:"class_id"`
	Title     string `form:"title" json:"title" binding:"required" label:"title"`
	Content   string `form:"content" json:"content" binding:"required" label:"content"`
	MdContent string `form:"md_content" json:"md_content" binding:"required" label:"md_content"`
}

type ArticleListRequest struct {
	Keyword  string `json:"keyword"`
	FindType int    `json:"find_type"`
	Cid      int    `json:"cid"`
	Page     int    `json:"page"`
}

type ArticleDetailRequest struct {
	ArticleId int `form:"article_id" json:"article_id"`
}

type ArticleClassEditRequest struct {
	ClassId   int    `form:"class_id" json:"class_id"`
	ClassName string `form:"class_name" json:"class_name" binding:"required" label:"class_name"`
}

type ArticleClassDeleteRequest struct {
	ClassId int `form:"class_id" json:"class_id" binding:"required"`
}

type ArticleClassSortRequest struct {
	ClassId  int `form:"class_id" json:"class_id" binding:"required"`
	SortType int `form:"sort_type" json:"sort_type" binding:"required,oneof=1 2"`
}

type ArticleTagEditRequest struct {
	TagId   int    `form:"tag_id" json:"tag_id"`
	TagName string `form:"tag_name" json:"tag_name" binding:"required" label:"tag_name"`
}

type ArticleTagDeleteRequest struct {
	TagId int `form:"tag_id" json:"tag_id" binding:"required"`
}
