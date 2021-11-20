package dao

import "context"

type TalkRecordsVoteDao struct {
}

func (dao *TalkRecordsVoteDao) GetVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	uids := make([]int, 0)

	return uids, nil
}

func (dao *TalkRecordsVoteDao) UpdateVoteAnswerUser(ctx context.Context, vid int) ([]int, error) {
	uids := make([]int, 0)

	return uids, nil
}

func (dao *TalkRecordsVoteDao) GetVoteStatistics(ctx context.Context, vid int) error {
	return nil
}

func (dao *TalkRecordsVoteDao) UpdateVoteStatistics(ctx context.Context, vid int) error {
	return nil
}
