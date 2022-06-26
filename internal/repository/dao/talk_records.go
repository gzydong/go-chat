package dao

import (
	"context"

	model2 "go-chat/internal/repository/model"
)

type TalkRecordsDao struct {
	*BaseDao
}

func NewTalkRecordsDao(baseDao *BaseDao) *TalkRecordsDao {
	return &TalkRecordsDao{BaseDao: baseDao}
}

// GetChatRecords 查询对话记录
func (dao *TalkRecordsDao) GetChatRecords() {

}

func (dao *TalkRecordsDao) SearchChatRecords() {

}

type FindFileRecordData struct {
	Record   *model2.TalkRecords
	FileInfo *model2.TalkRecordsFile
}

func (dao *TalkRecordsDao) FindFileRecord(ctx context.Context, recordId int) (*FindFileRecordData, error) {
	var (
		record   *model2.TalkRecords
		fileInfo *model2.TalkRecordsFile
	)

	if err := dao.Db().First(&record, recordId).Error; err != nil {
		return nil, err
	}

	if err := dao.Db().First(&fileInfo, "record_id = ?", recordId).Error; err != nil {
		return nil, err
	}

	return &FindFileRecordData{
		Record:   record,
		FileInfo: fileInfo,
	}, nil
}
