package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"path"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
	"go-chat/internal/repository/dao"
	"go-chat/internal/repository/model"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
)

type MultipartInitiateOpts struct {
	UserId int
	Name   string
	Size   int64
}

type MultipartUploadOpts struct {
	UserId     int
	UploadId   string
	SplitIndex int
	SplitNum   int
	File       *multipart.FileHeader
}

type SplitUploadService struct {
	*BaseService
	dao        *dao.SplitUploadDao
	conf       *config.Config
	fileSystem *filesystem.Filesystem
}

func NewSplitUploadService(baseService *BaseService, dao *dao.SplitUploadDao, conf *config.Config, fileSystem *filesystem.Filesystem) *SplitUploadService {
	return &SplitUploadService{BaseService: baseService, dao: dao, conf: conf, fileSystem: fileSystem}
}

func (s *SplitUploadService) Dao() *dao.SplitUploadDao {
	return s.dao
}

func (s *SplitUploadService) InitiateMultipartUpload(ctx context.Context, params *MultipartInitiateOpts) (*model.SplitUpload, error) {

	// 计算拆分数量 3M
	num := math.Ceil(float64(params.Size) / float64(3<<20))

	m := &model.SplitUpload{
		Type:         1,
		Drive:        entity.FileDriveMode(s.fileSystem.Driver()),
		UserId:       params.UserId,
		OriginalName: params.Name,
		SplitNum:     int(num),
		FileExt:      strings.TrimPrefix(path.Ext(params.Name), "."),
		FileSize:     params.Size,
		Path:         fmt.Sprintf("private/tmp/multipart/%s/%s.tmp", timeutil.DateNumber(), encrypt.Md5(strutil.Random(20))),
		Attr:         "{}",
	}

	uploadId, err := s.fileSystem.Default.InitiateMultipartUpload(m.Path, m.OriginalName)
	if err != nil {
		return nil, err
	}

	m.UploadId = uploadId

	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (s *SplitUploadService) MultipartUpload(ctx context.Context, opts *MultipartUploadOpts) error {
	info := &model.SplitUpload{}
	if err := s.Db().First(info, "upload_id = ? and type = 1", opts.UploadId).Error; err != nil {
		return err
	}

	stream, err := filesystem.ReadMultipartStream(opts.File)
	if err != nil {
		return err
	}

	dirPath := fmt.Sprintf("private/tmp/%s/%s/%d-%s.tmp", timeutil.DateNumber(), encrypt.Md5(opts.UploadId), opts.SplitIndex, opts.UploadId)

	data := &model.SplitUpload{
		Type:         2,
		Drive:        info.Drive,
		UserId:       opts.UserId,
		UploadId:     opts.UploadId,
		OriginalName: info.OriginalName,
		SplitIndex:   opts.SplitIndex,
		SplitNum:     opts.SplitNum,
		Path:         dirPath,
		FileExt:      info.FileExt,
		FileSize:     opts.File.Size,
		Attr:         "{}",
	}

	switch data.Drive {
	case entity.FileDriveLocal:
		_ = s.fileSystem.Default.Write(stream, data.Path)
	case entity.FileDriveCos:
		etag, err := s.fileSystem.Cos.UploadPart(info.Path, data.UploadId, data.SplitIndex+1, stream)
		if err != nil {
			return err
		}

		data.Attr = jsonutil.Encode(map[string]string{
			"etag": etag,
		})

	default:
		return errors.New("未知文件驱动类型")
	}

	if err := s.Db().Create(data).Error; err != nil {
		return err
	}

	// 判断是否为最后一个分片上传
	if opts.SplitNum == opts.SplitIndex+1 {
		err = s.merge(info)
	}

	return err
}

// combine
func (s *SplitUploadService) merge(info *model.SplitUpload) error {
	items, err := s.dao.GetSplitList(info.UploadId)
	if err != nil {
		return err
	}

	switch info.Drive {
	case entity.FileDriveLocal:
		for _, item := range items {
			stream, err := s.fileSystem.Default.ReadStream(item.Path)
			if err != nil {
				return err
			}

			if err := s.fileSystem.Local.AppendWrite(stream, info.Path); err != nil {
				return err
			}
		}
	case entity.FileDriveCos:
		opt := &cos.CompleteMultipartUploadOptions{}
		for _, item := range items {
			attr := make(map[string]string)

			if err := jsonutil.Decode(item.Attr, &attr); err != nil {
				return err
			}

			opt.Parts = append(opt.Parts, cos.Object{
				PartNumber: item.SplitIndex + 1,
				ETag:       attr["etag"],
			})
		}

		if err := s.fileSystem.Cos.CompleteMultipartUpload(info.Path, info.UploadId, opt); err != nil {
			return err
		}
	default:
		return errors.New("未知文件驱动类型")
	}

	return nil
}
