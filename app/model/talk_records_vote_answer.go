package model

type TalkRecordsVoteAnswer struct {
	ID        int    `json:"id" grom:"comment:自增ID"`
	VoteId    int    `json:"vote_id" grom:"comment:投票ID"`
	UserId    int    `json:"user_id" grom:"comment:投票用户"`
	Option    string `json:"option" grom:"comment:投票选项"`
	CreatedAt string `json:"created_at" grom:"comment:投票时间"`
}
