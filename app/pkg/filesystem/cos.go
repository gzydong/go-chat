package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go-chat/app/pkg/timeutil"
	"go-chat/config"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

// @link https://cloud.tencent.com/document/product/436/31215
type CosFilesystem struct {
	conf   *config.Config
	client *cos.Client
}

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

	return &CosFilesystem{conf: conf, client: client}
}

// Write 文件写入
func (c *CosFilesystem) Write(data []byte, filePath string) error {

	_, err := c.client.Object.Put(context.Background(), filePath, bytes.NewBuffer(data), &cos.ObjectPutOptions{})

	return err
}

// WriteLocal 本地文件上传
func (c *CosFilesystem) WriteLocal(localFile string, filePath string) error {
	_, _, err := c.client.Object.Upload(context.Background(), filePath, localFile, nil)

	return err
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
