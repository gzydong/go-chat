package service

import (
	"fmt"
	"go-chat/app/dao"
	"go-chat/app/model"
	"go-chat/app/pkg/slice"
	"strconv"
	"strings"
)

type EmoticonService struct {
	*BaseService
	Dao *dao.EmoticonDao
}

func NewEmoticonService(base *BaseService, dao *dao.EmoticonDao) *EmoticonService {
	return &EmoticonService{BaseService: base, Dao: dao}
}

func (s *EmoticonService) RemoveUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.Dao.GetUserInstallIds(uid)

	if !slice.InInt(emoticonId, ids) {
		return fmt.Errorf("数据不存在！")
	}

	items := make([]string, 0, len(ids)-1)

	for _, id := range ids {
		if id != emoticonId {
			items = append(items, strconv.Itoa(id))
		}
	}

	return s.db.Model(model.UsersEmoticon{}).Where("user_id = ?", uid).Update("emoticon_ids", strings.Join(items, ",")).Error
}

func (s *EmoticonService) AddUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.Dao.GetUserInstallIds(uid)

	if slice.InInt(emoticonId, ids) {
		return nil
	}

	ids = append(ids, emoticonId)

	fmt.Println(ids, slice.IntToIds(ids))

	return s.db.Model(model.UsersEmoticon{}).Where("user_id = ?", uid).Update("emoticon_ids", slice.IntToIds(ids)).Error
}
