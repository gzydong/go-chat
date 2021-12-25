package service

import (
	"bytes"
	"context"
	"fmt"
	"go-chat/app/dao"
	"go-chat/app/model"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/utils"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type EmoticonService struct {
	*BaseService
	dao        *dao.EmoticonDao
	fileSystem *filesystem.Filesystem
}

func NewEmoticonService(baseService *BaseService, dao *dao.EmoticonDao, fileSystem *filesystem.Filesystem) *EmoticonService {
	return &EmoticonService{BaseService: baseService, dao: dao, fileSystem: fileSystem}
}

func (s *EmoticonService) Dao() *dao.EmoticonDao {
	return s.dao
}

func (s *EmoticonService) RemoveUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.dao.GetUserInstallIds(uid)

	if !slice.InInt(emoticonId, ids) {
		return fmt.Errorf("数据不存在！")
	}

	items := make([]string, 0, len(ids)-1)

	for _, id := range ids {
		if id != emoticonId {
			items = append(items, strconv.Itoa(id))
		}
	}

	return s.db.Model(&model.UsersEmoticon{}).Where("user_id = ?", uid).Update("emoticon_ids", strings.Join(items, ",")).Error
}

func (s *EmoticonService) AddUserSysEmoticon(uid int, emoticonId int) error {
	ids := s.dao.GetUserInstallIds(uid)

	if slice.InInt(emoticonId, ids) {
		return nil
	}

	ids = append(ids, emoticonId)

	return s.db.Model(&model.UsersEmoticon{}).Where("user_id = ?", uid).Update("emoticon_ids", slice.IntToIds(ids)).Error
}

// DeleteCollect 删除自定义表情包
func (s *EmoticonService) DeleteCollect(uid int, ids []int) error {
	return s.db.Delete(&model.EmoticonItem{}, "id in ? and emoticon_id = 0 and user_id = ?", ids, uid).Error
}

func (s *EmoticonService) CustomizeUpload(ctx context.Context, uid int, file *multipart.FileHeader) (*model.EmoticonItem, error) {

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return nil, err
	}

	size := utils.ReadFileImage(bytes.NewReader(stream))
	ext := strutil.FileSuffix(file.Filename)
	src := fmt.Sprintf("public/media/image/emoticon/%s/%s", time.Now().Format("20060102"), strutil.GenImageName(ext, size.Width, size.Height))
	if err = s.fileSystem.Default.Write(stream, src); err != nil {
		return nil, err
	}

	m := &model.EmoticonItem{
		UserId:     uid,
		Describe:   "自定义表情包",
		Url:        s.fileSystem.Default.PublicUrl(src),
		FileSuffix: ext,
		FileSize:   int(file.Size),
	}

	if err := s.Db().Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}
