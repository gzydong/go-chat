package service

import (
	"fmt"
	"strconv"
	"strings"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type EmoticonService struct {
	*repo.Source
	emoticon   *repo.Emoticon
	filesystem *filesystem.Filesystem
}

func NewEmoticonService(baseService *repo.Source, repo *repo.Emoticon, fileSystem *filesystem.Filesystem) *EmoticonService {
	return &EmoticonService{Source: baseService, emoticon: repo, filesystem: fileSystem}
}

func (s *EmoticonService) Dao() *repo.Emoticon {
	return s.emoticon
}

func (s *EmoticonService) RemoveUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.emoticon.GetUserInstallIds(uid)

	if !sliceutil.Include(emoticonId, ids) {
		return fmt.Errorf("数据不存在！")
	}

	items := make([]string, 0, len(ids)-1)
	for _, id := range ids {
		if id != emoticonId {
			items = append(items, strconv.Itoa(id))
		}
	}

	return s.Db().Table("users_emoticon").Where("user_id = ?", uid).Update("emoticon_ids", strings.Join(items, ",")).Error
}

func (s *EmoticonService) AddUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.emoticon.GetUserInstallIds(uid)
	if sliceutil.Include(emoticonId, ids) {
		return nil
	}

	ids = append(ids, emoticonId)
	return s.Db().Table("users_emoticon").Where("user_id = ?", uid).Update("emoticon_ids", sliceutil.ToIds(ids)).Error
}

// DeleteCollect 删除自定义表情包
func (s *EmoticonService) DeleteCollect(uid int, ids []int) error {
	return s.Db().Delete(&model.EmoticonItem{}, "id in ? and emoticon_id = 0 and user_id = ?", ids, uid).Error
}
