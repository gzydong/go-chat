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
}

func NewQiniuFilesystem(conf *config.Config) *QiniuFilesystem {
	return &QiniuFilesystem{
		conf: conf,
	}
}

// Token 获取上传凭证
// todo token 需要加入缓存
func (s *QiniuFilesystem) Token() string {
	mac := qbox.NewMac(s.conf.Filesystem.Qiniu.AccessKey, s.conf.Filesystem.Qiniu.SecretKey)

	putPolicy := storage.PutPolicy{
		Scope:   s.conf.Filesystem.Qiniu.Bucket,
		Expires: 7200,
	}

	return putPolicy.UploadToken(mac)
}

func (s *QiniuFilesystem) Write(data []byte, filePath string) {
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
		fmt.Println(err)
		return
	}

	fmt.Println(ret.Key, ret.Hash)
}

func (s *QiniuFilesystem) WriteLocal(localFile string, filePath string) {
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
		fmt.Println(err)
		return
	}

	fmt.Println(ret.Key, ret.Hash)
}

func (s *QiniuFilesystem) Update() {
	fmt.Println("QiniuFilesystem :", "Update")
}

func (s *QiniuFilesystem) Rename() {

}

func (s *QiniuFilesystem) Copy() {

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
