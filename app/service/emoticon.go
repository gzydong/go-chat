package service

import "go-chat/app/dao"

type EmoticonService struct {
	*BaseService
	Dao *dao.EmoticonDao
}

func NewEmoticonService(dao *dao.EmoticonDao) *EmoticonService {
	return &EmoticonService{Dao: dao}
}
