package model

type ArticleTag struct {
	Id        int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 文章分类ID
	UserId    int    `gorm:"column:user_id;default:0" json:"user_id"`        // 用户ID
	TagName   string `gorm:"column:tag_name" json:"tag_name"`                // 标签名
	Sort      int    `gorm:"column:sort;default:0" json:"sort"`              // 排序
	CreatedAt int    `gorm:"column:created_at;default:0" json:"created_at"`  // 创建时间
}

type TagItem struct {
	Id      int    `json:"id"`       // 文章分类ID
	TagName string `json:"tag_name"` // 标签名
	Count   int    `json:"count"`    // 排序
}
