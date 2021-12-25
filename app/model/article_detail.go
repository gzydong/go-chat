package model

type ArticleDetail struct {
	Id        int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 文章详情ID
	ArticleId int    `gorm:"column:article_id" json:"article_id"`            // 文章ID
	MdContent string `gorm:"column:md_content" json:"md_content"`            // Markdown 内容
	Content   string `gorm:"column:content" json:"content"`                  // Markdown 解析HTML内容
}
