package repo

import (
	"context"

	"go-chat/internal/repository/model"
)

type TalkRecords struct {
	*Base
}

func NewTalkRecords(base *Base) *TalkRecords {
	return &TalkRecords{Base: base}
}

// GetChatRecords 查询对话记录
func (t *TalkRecords) GetChatRecords() {

}

func (t *TalkRecords) SearchChatRecords() {

}

type FindFileRecordData struct {
	Record   *model.TalkRecords
	FileInfo *model.TalkRecordsFile
}

func (t *TalkRecords) FindFileRecord(ctx context.Context, recordId int) (*FindFileRecordData, error) {
	var (
		record   *model.TalkRecords
		fileInfo *model.TalkRecordsFile
	)

	tx := t.Db.WithContext(ctx)

	if err := tx.First(&record, recordId).Error; err != nil {
		return nil, err
	}

	if err := tx.First(&fileInfo, "record_id = ?", recordId).Error; err != nil {
		return nil, err
	}

	return &FindFileRecordData{
		Record:   record,
		FileInfo: fileInfo,
	}, nil
}
