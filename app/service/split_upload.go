package service

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go-chat/app/dao"
	"go-chat/app/entity"
	"go-chat/app/http/request"
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

type InitiateParams struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	UserId int    `json:"user_id"`
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

// IsUploadFile 判断拆分文件上传ID是否存在
func (s *SplitUploadService) IsUploadFile(ctx context.Context, uid int, hashId string) {

}

func (s *SplitUploadService) InitiateMultipartUpload(ctx context.Context, params *InitiateParams) (*model.SplitUpload, error) {

	// 计算拆分数量
	num := math.Ceil(float64(params.Size) / float64(2<<20))

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

func (s *SplitUploadService) MultipartUpload(ctx context.Context, uid int, req *request.UploadMultipartRequest, file *multipart.FileHeader) (interface{}, error) {
	info := &model.SplitUpload{}
	if err := s.Db().First(info, "upload_id = ? and type = 1", req.UploadId).Error; err != nil {
		return nil, err
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return nil, err
	}

	data := &model.SplitUpload{
		Type:         2,
		Drive:        info.Drive,
		UserId:       uid,
		UploadId:     req.UploadId,
		OriginalName: req.Name,
		SplitIndex:   req.SplitIndex,
		SplitNum:     req.SplitNum,
		SaveDir:      fmt.Sprintf("private/tmp/%s/%s/%d-%s.tmp", timeutil.DateNumber(), req.UploadId, req.SplitIndex, req.UploadId),
		FileExt:      req.Ext,
		FileSize:     file.Size,
		Attr:         "{}",
	}

	switch data.Drive {
	case entity.FileDriveLocal:
		_ = s.fileSystem.Default.Write(stream, data.SaveDir)
	case entity.FileDriveCos:
		etag, err := s.fileSystem.Cos.UploadPart(info.SaveDir, data.UploadId, data.SplitIndex+1, stream)
		if err != nil {
			return nil, err
		}

		data.Attr = jsonutil.JsonEncode(map[string]string{
			"etag": etag,
		})
	}

	if err := s.Db().Create(data).Error; err != nil {
		return nil, err
	}

	// 判断是否为最后一个分片上传
	if req.SplitNum == req.SplitIndex+1 {
		_ = s.merge(info)
	}

	return nil, nil
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
