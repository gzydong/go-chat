package service

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

var _ IEmoticonService = (*EmoticonService)(nil)

type IEmoticonService interface {
	AddUserSysEmoticon(uid int, emoticonId int) error
	RemoveUserSysEmoticon(uid int, emoticonId int) error
	DeleteCollect(uid int, ids []int) error
}

type EmoticonService struct {
	*repo.Source
	EmoticonRepo *repo.Emoticon
	Filesystem   filesystem.IFilesystem
}

func (s *EmoticonService) RemoveUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.EmoticonRepo.GetUserInstallIds(uid)

	if !slices.Contains(ids, emoticonId) {
		return fmt.Errorf("数据不存在！")
	}

	items := make([]string, 0, len(ids)-1)
	for _, id := range ids {
		if id != emoticonId {
			items = append(items, strconv.Itoa(id))
		}
	}

	return s.Source.Db().Table("users_emoticon").Where("user_id = ?", uid).Update("emoticon_ids", strings.Join(items, ",")).Error
}

func (s *EmoticonService) AddUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.EmoticonRepo.GetUserInstallIds(uid)
	if slices.Contains(ids, emoticonId) {
		return nil
	}

	ids = append(ids, emoticonId)
	return s.Source.Db().Table("users_emoticon").Where("user_id = ?", uid).Update("emoticon_ids", sliceutil.ToIds(ids)).Error
}

// DeleteCollect 删除自定义表情包
func (s *EmoticonService) DeleteCollect(uid int, ids []int) error {
	return s.Source.Db().Delete(&model.EmoticonItem{}, "id in ? and emoticon_id = 0 and user_id = ?", ids, uid).Error
}
