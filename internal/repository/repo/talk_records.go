package repo

import (
	"context"

	model2 "go-chat/internal/repository/model"
)

type TalkRecords struct {
	*Base
}

func NewTalkRecords(baseDao *Base) *TalkRecords {
	return &TalkRecords{Base: baseDao}
}

// GetChatRecords 查询对话记录
func (repo *TalkRecords) GetChatRecords() {

}

func (repo *TalkRecords) SearchChatRecords() {

}

type FindFileRecordData struct {
	Record   *model2.TalkRecords
	FileInfo *model2.TalkRecordsFile
}

func (repo *TalkRecords) FindFileRecord(ctx context.Context, recordId int) (*FindFileRecordData, error) {
	var (
		record   *model2.TalkRecords
		fileInfo *model2.TalkRecordsFile
	)

	if err := repo.db.First(&record, recordId).Error; err != nil {
		return nil, err
	}

	if err := repo.db.First(&fileInfo, "record_id = ?", recordId).Error; err != nil {
		return nil, err
	}

	return &FindFileRecordData{
		Record:   record,
		FileInfo: fileInfo,
	}, nil
}
