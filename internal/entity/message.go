package entity

type TextMessage struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
	Mention struct {
		All  int      `json:"all"`
		Uids []string `json:"uids"`
	} `json:"mention"`
}

type ImageMessage struct {
	Type   int    `json:"type"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type VoiceMessage struct {
	Type     int    `json:"type"`
	Url      string `json:"url"`
	TimeTrad int    `json:"timeTrad"`
}

type FileMessage struct {
	Type     int    `json:"type"`
	UploadId string `json:"upload_id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
}

// SysCreateGroupMessage 创建群聊
type SysCreateGroupMessage struct {
	Type        int    `json:"type"`
	Creator     string `json:"creator"`
	CreatorName string `json:"creator_name"`
	Content     string `json:"content"`
	Extra       []struct {
		Uid  string `json:"uid"`
		Name string `json:"name"`
	} `json:"extra"`
}
