package model

type Emoticon struct {
	ID        int      `json:"id" grom:"comment:分组ID"`
	Name      string   `json:"name" grom:"comment:分组名称"`
	Icon      string   `json:"icon" grom:"comment:分组图标"`
	Status    int      `json:"status" grom:"comment:分组状态"`
	CreatedAt int64    `json:"created_at" grom:"comment:创建时间"`
	UpdatedAt DateTime `json:"updated_at" grom:"comment:更新时间"`
}
