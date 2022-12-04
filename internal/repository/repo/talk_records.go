package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecords struct {
	ichat.Repo[model.TalkRecords]
}

func NewTalkRecords(db *gorm.DB) *TalkRecords {
	return &TalkRecords{Repo: ichat.NewRepo[model.TalkRecords](db)}
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
