package api

type SysEmoticonList struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	IsUsed int    `json:"status"`
}
