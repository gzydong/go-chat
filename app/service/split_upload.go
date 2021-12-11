package service

import (
	"context"
	"go-chat/app/model"
	"go-chat/app/pkg/strutil"
	"math"
	"path"
	"strings"
)

type SplitUploadService struct {
	*BaseService
}

func NewSplitUploadService(baseService *BaseService) *SplitUploadService {
	return &SplitUploadService{BaseService: baseService}
}

type InitiateParams struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	UserId int    `json:"user_id"`
}

func (s *SplitUploadService) InitiateMultipartUpload(ctx context.Context, params *InitiateParams) (*model.FileSplitUpload, error) {

	// 计算拆分数量
	num := math.Ceil(float64(params.Size) / float64(2<<20))

	m := &model.FileSplitUpload{
		FileType:     1,
		UserId:       params.UserId,
		HashName:     strutil.Md5([]byte(strutil.GenRandomString(20))),
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

func (s *SplitUploadService) MultipartAppendUpload(ctx context.Context) {

}
