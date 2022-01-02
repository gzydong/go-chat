package service

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go-chat/app/dao"
	"go-chat/app/entity"
	"go-chat/app/model"
	"go-chat/app/pkg/encrypt"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/jsonutil"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"go-chat/config"
	"math"
	"mime/multipart"
	"path"
	"strings"
)

type MultipartInitiateOpts struct {
	UserId int
	Name   string
	Size   int64
}

type MultipartUploadOpts struct {
	UserId     int
	UploadId   string
	Name       string
	Ext        string
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
		SaveDir:      fmt.Sprintf("private/tmp/multipart/%s/%s.tmp", timeutil.DateNumber(), encrypt.Md5(strutil.Random(20))),
		Attr:         "{}",
	}

	uploadId, err := s.fileSystem.Default.InitiateMultipartUpload(m.SaveDir, m.OriginalName)
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

	data := &model.SplitUpload{
		Type:         2,
		Drive:        info.Drive,
		UserId:       opts.UserId,
		UploadId:     opts.UploadId,
		OriginalName: opts.Name,
		SplitIndex:   opts.SplitIndex,
		SplitNum:     opts.SplitNum,
		SaveDir:      fmt.Sprintf("private/tmp/%s/%s/%d-%s.tmp", timeutil.DateNumber(), opts.UploadId, opts.SplitIndex, opts.UploadId),
		FileExt:      opts.Ext,
		FileSize:     opts.File.Size,
		Attr:         "{}",
	}

	switch data.Drive {
	case entity.FileDriveLocal:
		_ = s.fileSystem.Default.Write(stream, data.SaveDir)
	case entity.FileDriveCos:
		etag, err := s.fileSystem.Cos.UploadPart(info.SaveDir, data.UploadId, data.SplitIndex+1, stream)
		if err != nil {
			return err
		}

		data.Attr = jsonutil.JsonEncode(map[string]string{
			"etag": etag,
		})
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
			stream, err := s.fileSystem.Default.ReadStream(item.SaveDir)
			if err != nil {
				fmt.Println("ReadContent err:", err.Error())
				return err
			}

			if err := s.fileSystem.Local.AppendWrite(stream, info.SaveDir); err != nil {
				fmt.Println("AppendWrite err:", err)
				return err
			}
		}
	case entity.FileDriveCos:
		opt := &cos.CompleteMultipartUploadOptions{}
		for _, item := range items {
			attr := make(map[string]string)

			if err := jsonutil.JsonDecode(item.Attr, &attr); err != nil {
				return err
			}

			opt.Parts = append(opt.Parts, cos.Object{
				PartNumber: item.SplitIndex + 1,
				ETag:       attr["etag"],
			})
		}

		if err := s.fileSystem.Cos.CompleteMultipartUpload(info.SaveDir, info.UploadId, opt); err != nil {
			return err
		}
	}

	return nil
}
