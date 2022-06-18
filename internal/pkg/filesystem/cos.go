package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	"go-chat/config"
	"go-chat/internal/pkg/timeutil"
)

type CosFilesystem struct {
	conf   *config.Config
	client *cos.Client
}

// NewCosFilesystem ...
// See: https://cloud.tencent.com/document/product/436/31215
func NewCosFilesystem(conf *config.Config) *CosFilesystem {

	bucketURL, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", conf.Filesystem.Cos.Bucket, conf.Filesystem.Cos.Region))

	client := cos.NewClient(&cos.BaseURL{
		BucketURL: bucketURL,
	}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.Filesystem.Cos.SecretId,
			SecretKey: conf.Filesystem.Cos.SecretKey,
		},
	})

	return &CosFilesystem{conf, client}
}

// Write 文件写入
func (c *CosFilesystem) Write(data []byte, filePath string) error {

	_, err := c.client.Object.Put(context.Background(), filePath, bytes.NewBuffer(data), nil)

	return err
}

// WriteLocal 本地文件上传
func (c *CosFilesystem) WriteLocal(localFile string, filePath string) error {
	_, _, err := c.client.Object.Upload(context.Background(), filePath, localFile, nil)

	return err
}

func (c *CosFilesystem) WriteFromFile(file *multipart.FileHeader, filePath string) error {
	return nil
}

// Copy 文件拷贝
func (c *CosFilesystem) Copy(srcPath, filePath string) error {

	sourceURL := fmt.Sprintf("%s/%s", c.client.BaseURL.BucketURL.Host, srcPath)

	_, _, err := c.client.Object.Copy(context.Background(), filePath, sourceURL, nil)

	return err
}

// Delete 删除一个文件或空文件夹
func (c *CosFilesystem) Delete(filePath string) error {

	_, err := c.client.Object.Delete(context.Background(), filePath, nil)

	return err
}

// DeleteDir 删除文件夹
func (c *CosFilesystem) DeleteDir(path string) error {
	path = strings.TrimSuffix(path, "/") + "/"

	return c.Delete(path)
}

// CreateDir 递归创建文件夹
func (c *CosFilesystem) CreateDir(path string) error {

	path = strings.TrimSuffix(path, "/") + "/"

	_, err := c.client.Object.Put(context.Background(), path, strings.NewReader(""), nil)

	return err
}

// Stat 文件信息
func (c *CosFilesystem) Stat(filePath string) (*FileStat, error) {
	resp, err := c.client.Object.Head(context.Background(), filePath, nil)
	if err != nil {
		return nil, err
	}

	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	t, _ := time.ParseInLocation(time.RFC1123, resp.Header.Get("Last-Modified"), timeutil.Location())

	return &FileStat{
		Name:        path.Base(filePath),
		Size:        int64(size),
		Ext:         path.Ext(filePath),
		LastModTime: t.Add(8 * time.Hour),
		MimeType:    resp.Header.Get("Content-Type"),
	}, nil
}

func (c *CosFilesystem) Append(filePath string) {

}

func (c *CosFilesystem) PublicUrl(filePath string) string {
	return c.client.Object.GetObjectURL(filePath).String()
}

func (c *CosFilesystem) PrivateUrl(filePath string, timeout int) string {
	result, err := c.client.Object.GetPresignedURL(
		context.Background(),
		http.MethodGet, filePath,
		c.conf.Filesystem.Cos.SecretId,
		c.conf.Filesystem.Cos.SecretKey,
		time.Second*time.Duration(timeout),
		nil,
	)

	if err != nil {
		return ""
	}

	return result.String()
}

// ReadStream 读取文件流信息
func (c *CosFilesystem) ReadStream(filePath string) ([]byte, error) {
	return nil, nil
}

func (c *CosFilesystem) InitiateMultipartUpload(filePath string, fileName string) (string, error) {
	opt := &cos.InitiateMultipartUploadOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentDisposition: fmt.Sprintf("attachment; filename=\"%s\"", fileName),
		},
	}

	resp, _, err := c.client.Object.InitiateMultipartUpload(context.Background(), filePath, opt)
	if err != nil {
		return "", err
	}

	return resp.UploadID, nil
}

func (c *CosFilesystem) UploadPart(filePath string, uploadID string, num int, stream []byte) (string, error) {
	resp, err := c.client.Object.UploadPart(context.Background(), filePath, uploadID, num, bytes.NewBuffer(stream), nil)

	if err != nil {
		return "", err
	}

	return strings.Trim(resp.Header.Get("ETag"), `"`), nil
}

func (c *CosFilesystem) CompleteMultipartUpload(filePath string, uploadID string, opt *cos.CompleteMultipartUploadOptions) error {
	_, _, err := c.client.Object.CompleteMultipartUpload(context.Background(), filePath, uploadID, opt)

	return err
}
