package model

type TalkRecordsVoteAnswer struct {
	ID        int    `json:"id"`
	VoteId    int    `json:"vote_id"`
	UserId    int    `json:"user_id"`
	Option    string `json:"option"`
	CreatedAt string `json:"created_at"`
}
