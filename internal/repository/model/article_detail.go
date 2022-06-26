package model

type ArticleDetail struct {
	Id        int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`         // 文章详情ID
	ArticleId int    `gorm:"column:article_id;default:0;NOT NULL" json:"article_id"` // 文章ID
	MdContent string `gorm:"column:md_content;NOT NULL" json:"md_content"`           // Markdown 内容
	Content   string `gorm:"column:content;NOT NULL" json:"content"`                 // Markdown 解析HTML内容
}
