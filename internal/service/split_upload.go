package service

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
)

var _ ISplitUploadService = (*SplitUploadService)(nil)

type ISplitUploadService interface {
	InitiateMultipartUpload(ctx context.Context, params *MultipartInitiateOpt) (*model.SplitUpload, error)
	MultipartUpload(ctx context.Context, opt *MultipartUploadOpt) error
}

type SplitUploadService struct {
	*repo.Source
	SplitUploadRepo *repo.SplitUpload
	Config          *config.Config
	FileSystem      filesystem.IFilesystem
}

type MultipartInitiateOpt struct {
	UserId int
	Name   string
	Size   int64
}

func (s *SplitUploadService) InitiateMultipartUpload(ctx context.Context, params *MultipartInitiateOpt) (*model.SplitUpload, error) {
	// 计算拆分数量 5M
	num := math.Ceil(float64(params.Size) / float64(5*1024*1024))

	now := time.Now()
	m := &model.SplitUpload{
		Type:         1,
		Drive:        entity.FileDriveMode(s.FileSystem.Driver()),
		UserId:       params.UserId,
		OriginalName: params.Name,
		SplitNum:     int(num),
		FileExt:      strings.TrimPrefix(path.Ext(params.Name), "."),
		FileSize:     params.Size,
		Path:         fmt.Sprintf("multipart/%s/%s.tmp", now.Format("20060102"), uuid.New().String()),
		Attr:         "{}",
	}

	uploadId, err := s.FileSystem.InitiateMultipartUpload(s.FileSystem.BucketPrivateName(), m.Path)
	if err != nil {
		return nil, err
	}

	m.UploadId = uploadId

	if err := s.Source.Db().WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

type MultipartUploadOpt struct {
	UserId     int
	UploadId   string
	SplitIndex int
	SplitNum   int
	File       *multipart.FileHeader
}

func (s *SplitUploadService) MultipartUpload(ctx context.Context, opt *MultipartUploadOpt) error {
	info, err := s.SplitUploadRepo.FindByWhere(ctx, "upload_id = ? and type = 1", opt.UploadId)
	if err != nil {
		return err
	}

	stream, err := filesystem.ReadMultipartStream(opt.File)
	if err != nil {
		return err
	}

	data := &model.SplitUpload{
		Type:         2,
		Drive:        info.Drive,
		UserId:       opt.UserId,
		UploadId:     opt.UploadId,
		OriginalName: info.OriginalName,
		SplitIndex:   opt.SplitIndex,
		SplitNum:     opt.SplitNum,
		Path:         "",
		FileExt:      info.FileExt,
		FileSize:     opt.File.Size,
		Attr:         "{}",
	}

	read := bytes.NewReader(stream)

	objectPart, err := s.FileSystem.PutObjectPart(
		s.FileSystem.BucketPrivateName(),
		info.Path,
		info.UploadId,
		opt.SplitIndex,
		read,
		read.Size(),
	)
	if err != nil {
		return err
	}

	if objectPart.PartObjectName != "" {
		data.Path = objectPart.PartObjectName
	}

	data.Attr = jsonutil.Encode(objectPart)

	if err = s.Source.Db().Create(data).Error; err != nil {
		fmt.Println("ERR====>", err)
		return err
	}

	// 判断是否为最后一个分片上传
	if opt.SplitNum == opt.SplitIndex {
		err = s.merge(info)
	}

	return err
}

// combine
func (s *SplitUploadService) merge(info *model.SplitUpload) error {
	items, err := s.SplitUploadRepo.FindAll(context.Background(), func(db *gorm.DB) {
		db.Where("upload_id =? and type = 2", info.UploadId).Order("split_index asc")
	})

	if err != nil {
		return err
	}

	parts := make([]filesystem.ObjectPart, 0)
	for _, item := range items {
		var obj filesystem.ObjectPart
		if err = jsonutil.Decode(item.Attr, &obj); err != nil {
			return err
		}

		parts = append(parts, obj)
	}

	// 合并文件
	if err := s.FileSystem.CompleteMultipartUpload(s.FileSystem.BucketPrivateName(), info.Path, info.UploadId, parts); err != nil {
		return err
	}

	return nil
}
