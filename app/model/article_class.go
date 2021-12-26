package model

type ArticleClass struct {
	Id        int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 文章分类ID
	UserId    int    `gorm:"column:user_id;default:0" json:"user_id"`        // 用户ID
	ClassName string `gorm:"column:class_name" json:"class_name"`            // 分类名
	Sort      int    `gorm:"column:sort;default:0" json:"sort"`              // 排序
	IsDefault int    `gorm:"column:is_default;default:0" json:"is_default"`  // 默认分类1:是 0:不是
	CreatedAt int    `gorm:"column:created_at;default:0" json:"created_at"`  // 创建时间
}

type ArticleClassItem struct {
	Id        int    `json:"id"`         // 文章分类ID
	ClassName string `json:"class_name"` // 分类名
	IsDefault int    `json:"is_default"` // 默认分类1:是 0:不是
	Count     int    `json:"count"`      // 分类名
}
