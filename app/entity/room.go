package entity

type RoomType string

const (
	RoomGroupChat     RoomType = "room_group_chat"     // 群聊房间
	RoomLiveBroadcast RoomType = "room_live_broadcast" // 直播房间
)
