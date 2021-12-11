package filesystem

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go-chat/config"
	"os"
	"strings"
)

type OssFilesystem struct {
	conf   *config.Config
	client *oss.Client
	bucket *oss.Bucket
}

func NewOssFilesystem(conf *config.Config) *OssFilesystem {
	client, err := oss.New(
		conf.Filesystem.Oss.Endpoint,
		conf.Filesystem.Oss.AccessID,
		conf.Filesystem.Oss.AccessSecret,
		oss.EnableCRC(true),
	)

	if err != nil {
		fmt.Println("Error:", err)
	}

	// 获取存储空间。
	bucket, err := client.Bucket(conf.Filesystem.Oss.Bucket)

	if err != nil {
		panic(err)
	}

	return &OssFilesystem{
		conf:   conf,
		client: client,
		bucket: bucket,
	}
}

func (s *OssFilesystem) Write(data []byte, filePath string) error {
	return s.bucket.PutObject(filePath, bytes.NewReader(data))
}

func (s *OssFilesystem) WriteLocal(localFile string, filePath string) error {
	fd, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer fd.Close()

	return s.bucket.PutObject(filePath, fd)
}

// Copy 拷贝文件到同一个存储空间的另一个文件。
func (s *OssFilesystem) Copy(srcPath, filePath string) error {
	_, err := s.bucket.CopyObject(srcPath, filePath)

	return err
}

// Delete 删除文件
func (s *OssFilesystem) Delete(filePath string) error {
	return s.bucket.DeleteObject(filePath)
}

func (s *OssFilesystem) DeleteDir(path string) error {
	return nil
}

func (s *OssFilesystem) CreateDir(path string) error {
	return nil
}

// IsObjectExist 判断文件是否存在
func (s *OssFilesystem) IsObjectExist(filePath string) bool {
	isExist, _ := s.bucket.IsObjectExist(filePath)

	return isExist
}

func (s *OssFilesystem) Stat(filePath string) (*FileStat, error) {
	// 获取文件元信息。
	props, err := s.bucket.GetObjectMeta(filePath)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Object Meta: %#v", props)

	return &FileStat{
		//LastModTime: props.Get("X-Oss-server-Time"),
	}, nil
}

func (s *OssFilesystem) PublicUrl(filePath string) string {
	return fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(s.conf.Filesystem.Oss.Endpoint, "/"),
		strings.Trim(filePath, "/"),
	)
}

func (s *OssFilesystem) PrivateUrl(filePath string, timeout int) string {
	return ""
}
