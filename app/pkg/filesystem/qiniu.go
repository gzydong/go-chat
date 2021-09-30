package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	_ "github.com/qiniu/go-sdk/v7/storage"
	"go-chat/config"
	"strings"
)

// QiniuFilesystem
// @link 对接文档 https://developer.qiniu.com/kodo/1238/go#upload-flow
type QiniuFilesystem struct {
	conf *config.Config
	mac  *qbox.Mac
}

func NewQiniuFilesystem(conf *config.Config) *QiniuFilesystem {
	return &QiniuFilesystem{
		conf: conf,
		mac:  qbox.NewMac(conf.Filesystem.Qiniu.AccessKey, conf.Filesystem.Qiniu.SecretKey),
	}
}

// Token 获取上传凭证
// todo token 需要加入缓存
func (s *QiniuFilesystem) Token() string {
	putPolicy := storage.PutPolicy{
		Scope:   s.conf.Filesystem.Qiniu.Bucket,
		Expires: 7200,
	}

	return putPolicy.UploadToken(s.mac)
}

func (s *QiniuFilesystem) Write(data []byte, filePath string) error {
	filePath = strings.TrimLeft(filePath, "/")

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS:      true,                 // 是否使用https域名
		UseCdnDomains: false,                // 上传是否使用CDN上传加速
	}

	// 七牛标准的上传回复内容
	ret := storage.PutRet{}

	// 可选配置
	params := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)

	err := formUploader.Put(context.Background(), &ret, s.Token(), filePath, bytes.NewReader(data), int64(len(data)), &params)
	if err != nil {
		return err
	}

	fmt.Println(ret.Key, ret.Hash)

	return nil
}

func (s *QiniuFilesystem) WriteLocal(localFile string, filePath string) error {
	filePath = strings.TrimLeft(filePath, "/")

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS:      true,                 // 是否使用https域名
		UseCdnDomains: false,                // 上传是否使用CDN上传加速
	}

	// 七牛标准的上传回复内容
	ret := storage.PutRet{}

	// 可选配置
	params := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	err := formUploader.PutFile(context.Background(), &ret, s.Token(), filePath, localFile, &params)
	if err != nil {
		return err
	}

	return nil
}

func (s *QiniuFilesystem) Copy(srcPath, filePath string) error {
	cfg := storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuadong, // 空间对应的机房
	}

	bucketManager := storage.NewBucketManager(s.mac, &cfg)

	bucket := s.conf.Filesystem.Qiniu.Bucket

	return bucketManager.Copy(bucket, srcPath, bucket, filePath, false)
}

func (s *QiniuFilesystem) Delete(filePath string) error {
	return nil
}

func (s *QiniuFilesystem) DeleteDir(path string) error {
	return nil
}

func (s *QiniuFilesystem) CreateDir(path string) error {
	return nil
}

func (s *QiniuFilesystem) Stat(filePath string) {
	cfg := storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuadong, // 空间对应的机房
	}

	bucketManager := storage.NewBucketManager(s.mac, &cfg)

	bucket := s.conf.Filesystem.Qiniu.Bucket

	fileInfo, sErr := bucketManager.Stat(bucket, filePath)
	if sErr != nil {
		fmt.Println(sErr)
		return
	}

	fmt.Println(fileInfo.String())
	//可以解析文件的PutTime
	fmt.Println(storage.ParsePutTime(fileInfo.PutTime))
}
