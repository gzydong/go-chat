package web

type SetSystemEmoticonRequest struct {
	EmoticonId int `form:"emoticon_id" json:"emoticon_id" binding:"required,numeric" label:"emoticon_id"`
	Type       int `form:"type" json:"type" binding:"required,oneof=1 2" label:"type"`
}

type DeleteCollectRequest struct {
	Ids string `form:"ids" json:"ids" binding:"required,ids" label:"ids"`
}

type SysEmoticonList struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Status int    `json:"status"`
}

type EmoticonItem struct {
	MediaId int    `json:"media_id"`
	Src     string `json:"src"`
}

type SysEmoticonResponse struct {
	EmoticonId int             `json:"emoticon_id"`
	Url        string          `json:"url"`
	Name       string          `json:"name"`
	List       []*EmoticonItem `json:"list"`
}

type EmojiGroup struct {
	Label    string                `json:"label"`
	Icon     string                `json:"icon"`
	Children []*EmojiGroupChildren `json:"children"`
}

type EmojiGroupChildren struct {
	MediaId int    `json:"media_id"`
	Src     string `json:"src"`
}
