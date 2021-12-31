package request

type ArticleEditRequest struct {
	ArticleId int    `form:"article_id" json:"article_id"`
	ClassId   int    `form:"class_id" json:"class_id"`
	Title     string `form:"title" json:"title" binding:"required" label:"title"`
	Content   string `form:"content" json:"content" binding:"required" label:"content"`
	MdContent string `form:"md_content" json:"md_content" binding:"required" label:"md_content"`
}

type ArticleListRequest struct {
	Keyword  string `form:"keyword" json:"keyword"`
	FindType int    `form:"find_type" json:"find_type"`
	Cid      int    `form:"cid" json:"cid"`
	Page     int    `form:"page" json:"page"`
}

type ArticleDetailRequest struct {
	ArticleId int `form:"article_id" json:"article_id"`
}

type ArticleAsteriskRequest struct {
	ArticleId int `form:"article_id" json:"article_id" binding:"required"`
	Type      int `form:"article_id" json:"type" binding:"required,oneof=1 2"`
}

type ArticleTagsRequest struct {
	ArticleId int   `form:"article_id" json:"article_id" binding:"required"`
	Tags      []int `form:"tags" json:"tags"`
}

type ArticleMoveRequest struct {
	ArticleId int `form:"article_id" json:"article_id" binding:"required,gt=0"`
	ClassId   int `form:"class_id" json:"class_id" binding:"required,gt=0"`
}

type ArticleDeleteRequest struct {
	ArticleId int `form:"article_id" json:"article_id" binding:"required"`
}

type ArticleRecoverRequest struct {
	ArticleId int `form:"article_id" json:"article_id" binding:"required"`
}

type ArticleForeverDeleteRequest struct {
	ArticleId int `form:"article_id" json:"article_id" binding:"required"`
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

type ArticleAnnexUploadRequest struct {
	ArticleId int `form:"article_id" json:"article_id" binding:"required"`
}

type ArticleAnnexDownloadRequest struct {
	AnnexId int `form:"annex_id" json:"annex_id" binding:"required"`
}

type ArticleAnnexDeleteRequest struct {
	AnnexId int `form:"annex_id" json:"annex_id" binding:"required"`
}

type ArticleAnnexRecoverRequest struct {
	AnnexId int `form:"annex_id" json:"annex_id" binding:"required"`
}

type ArticleAnnexForeverDeleteRequest struct {
	AnnexId int `form:"annex_id" json:"annex_id" binding:"required"`
}
