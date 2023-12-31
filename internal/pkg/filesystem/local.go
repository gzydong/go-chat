package filesystem

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var _ IFilesystem = (*LocalFilesystem)(nil)

type LocalFilesystem struct {
	config LocalSystemConfig
}

func NewLocalFilesystem(config LocalSystemConfig) *LocalFilesystem {
	return &LocalFilesystem{config}
}

func (l LocalFilesystem) Driver() string {
	return LocalDriver
}

func (l LocalFilesystem) BucketPublicName() string {
	return l.config.BucketPublic
}

func (l LocalFilesystem) BucketPrivateName() string {
	return l.config.BucketPrivate
}

func (l LocalFilesystem) Stat(bucketName string, objectName string) (*FileStatInfo, error) {
	info, err := os.Stat(l.Path(bucketName, objectName))
	if err != nil {
		return nil, err
	}

	return &FileStatInfo{
		Name:        filepath.Base(objectName),
		Size:        info.Size(),
		Ext:         filepath.Ext(objectName),
		MimeType:    "",
		LastModTime: info.ModTime(),
	}, nil
}

func (l LocalFilesystem) Write(bucketName string, objectName string, stream []byte) error {
	filePath := l.Path(bucketName, objectName)

	dir := path.Dir(filePath)
	if len(dir) > 0 && !isDirExist(dir) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(stream)
	return err
}

// WriteLocal 本地文件上传
func (l LocalFilesystem) WriteLocal(bucketName string, localFile string, objectName string) error {
	srcFile, err := os.Open(localFile)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	objectName = l.Path(bucketName, objectName)
	dir := path.Dir(objectName)
	if len(dir) > 0 && !isDirExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	dstFile, err := os.OpenFile(objectName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (l LocalFilesystem) Copy(bucketName string, srcObjectName, objectName string) error {
	return l.WriteLocal(bucketName, l.Path(bucketName, srcObjectName), objectName)
}

func (l LocalFilesystem) CopyObject(srcBucketName string, srcObjectName, dstBucketName string, dstObjectName string) error {
	srcFile, err := os.Open(l.Path(srcBucketName, srcObjectName))
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstObjectName = l.Path(dstBucketName, dstObjectName)

	if dir := path.Dir(dstObjectName); len(dir) > 0 && !isDirExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	dstFile, err := os.OpenFile(dstObjectName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (l LocalFilesystem) Delete(bucketName string, objectName string) error {
	return os.Remove(l.Path(bucketName, objectName))
}

func (l LocalFilesystem) GetObject(bucketName string, objectName string) ([]byte, error) {
	return os.ReadFile(l.Path(bucketName, objectName))
}

func (l LocalFilesystem) PublicUrl(bucketName, objectName string) string {
	domain := fmt.Sprintf("http://%s", l.config.Endpoint)
	if l.config.SSL {
		domain = fmt.Sprintf("https://%s", l.config.Endpoint)
	}

	return fmt.Sprintf(
		"%s/%s/%s",
		strings.TrimRight(domain, "/"),
		bucketName,
		strings.Trim(objectName, "/"),
	)
}

func (l LocalFilesystem) PrivateUrl(bucketName, objectName string, _ string, _ time.Duration) string {
	return l.PublicUrl(bucketName, objectName)
}

func (l LocalFilesystem) InitiateMultipartUpload(_, _ string) (string, error) {
	return uuid.New().String(), nil
}

func (l LocalFilesystem) PutObjectPart(bucketName, _ string, uploadID string, index int, data io.Reader, _ int64) (ObjectPart, error) {
	stream, _ := io.ReadAll(data)

	objectName := fmt.Sprintf("multipart/%s/%d_%s.tmp", uploadID, index, uploadID)
	if err := l.Write(bucketName, objectName, stream); err != nil {
		return ObjectPart{}, err
	}

	return ObjectPart{
		ETag:           "",
		PartNumber:     index,
		PartObjectName: objectName,
	}, nil
}

func (l LocalFilesystem) CompleteMultipartUpload(bucketName, objectName, _ string, parts []ObjectPart) error {
	for _, part := range parts {
		stream, err := l.GetObject(bucketName, part.PartObjectName)
		if err != nil {
			return err
		}

		if err := l.appendWrite(bucketName, objectName, stream); err != nil {
			return err
		}
	}

	return nil
}

func (l LocalFilesystem) AbortMultipartUpload(bucketName, objectName, uploadID string) error {
	// TODO implement me
	panic("implement me")
}

func (l LocalFilesystem) appendWrite(bucketName, objectName string, stream []byte) error {
	filePath := l.Path(bucketName, objectName)

	dir := path.Dir(filePath)
	if len(dir) > 0 && !isDirExist(dir) {
		if err := os.MkdirAll(dir, 0766); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0766)
	if err != nil {
		return err
	}

	_, err = f.Write(stream)
	return err
}

// Path 获取文件地址绝对路径
func (l LocalFilesystem) Path(bucketName string, objectName string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		strings.TrimRight(l.config.Root, "/"),
		bucketName,
		strings.TrimLeft(objectName, "/"),
	)
}
