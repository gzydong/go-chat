package service

import (
	"context"
	"fmt"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/encrypt"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"io/ioutil"
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
	fileSystem *filesystem.Filesystem
}

func NewSplitUploadService(baseService *BaseService, fileSystem *filesystem.Filesystem) *SplitUploadService {
	return &SplitUploadService{BaseService: baseService, fileSystem: fileSystem}
}

// IsUploadFile 判断拆分文件上传ID是否存在
func (s *SplitUploadService) IsUploadFile(ctx context.Context, uid int, hashId string) {

}

func (s *SplitUploadService) InitiateMultipartUpload(ctx context.Context, params *InitiateParams) (*model.FileSplitUpload, error) {

	// 计算拆分数量
	num := math.Ceil(float64(params.Size) / float64(2<<20))

	m := &model.FileSplitUpload{
		FileType:     1,
		UserId:       params.UserId,
		HashName:     encrypt.Md5(strutil.Random(20)),
		OriginalName: params.Name,
		SplitNum:     int(num),
		FileExt:      strings.TrimPrefix(path.Ext(params.Name), "."),
		FileSize:     params.Size,
	}

	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (s *SplitUploadService) MultipartAppendUpload(ctx context.Context, req *request.UploadMultipartRequest, file *multipart.FileHeader) (interface{}, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer src.Close()

	content, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	_ = s.fileSystem.Write(content, fmt.Sprintf("private/tmp/%s/%s/%d-%s.tmp", timeutil.DateDay(), req.Hash, req.SplitIndex, req.Hash))

	return nil, nil
}
