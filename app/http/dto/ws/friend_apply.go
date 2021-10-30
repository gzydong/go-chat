package ws

type FriendApplyMessageContent struct {
	SenderID   int                       `json:"sender_id"`
	ReceiverID int                       `json:"receiver_id"`
	Remark     string                    `json:"remark"`
	Friend     *FriendApplyMessageFriend `json:"friend"`
}

type FriendApplyMessageFriend struct {
	UserID   int    `json:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
}

type FriendApplyCallbackMessageContent struct {
	SenderID   int                       `json:"sender_id"`
	ReceiverID int                       `json:"receiver_id"`
	Status     int                       `json:"status"`
	Remark     string                    `json:"remark"`
	Friend     *FriendApplyMessageFriend `json:"friend"`
}
