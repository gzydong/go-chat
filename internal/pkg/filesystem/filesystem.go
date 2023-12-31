package filesystem

import (
	"io"
	"time"
)

const (
	MinioDriver = "minio"
	LocalDriver = "local"
)

type IFilesystem interface {
	// Driver 驱动方式
	Driver() string

	// BucketPublicName 获取公开桶名
	BucketPublicName() string

	// BucketPrivateName 获取私有桶名
	BucketPrivateName() string

	// Stat 文件信息
	Stat(bucketName string, objectName string) (*FileStatInfo, error)

	// Write 文件写入
	Write(bucketName string, objectName string, stream []byte) error

	// Copy 文件拷贝
	Copy(bucketName string, srcObjectName, objectName string) error

	// CopyObject 文件拷贝
	CopyObject(srcBucketName string, srcObjectName, dstBucketName string, dstObjectName string) error

	// Delete 删除文件
	Delete(bucketName string, objectName string) error

	// GetObject 读取文件内容
	GetObject(bucketName string, objectName string) ([]byte, error)

	// PublicUrl 获取公开文件的访问地址
	PublicUrl(bucketName, objectName string) string

	// PrivateUrl 获取私有文件的访问地址
	PrivateUrl(bucketName, objectName string, filename string, expire time.Duration) string

	// InitiateMultipartUpload 初始化分片上传
	InitiateMultipartUpload(bucketName, objectName string) (string, error)

	// PutObjectPart 分片上传
	PutObjectPart(bucketName, objectName string, uploadID string, index int, data io.Reader, size int64) (ObjectPart, error)

	// CompleteMultipartUpload 完成分片上传
	CompleteMultipartUpload(bucketName, objectName, uploadID string, parts []ObjectPart) error

	// AbortMultipartUpload 取消分片上传
	AbortMultipartUpload(bucketName, objectName, uploadID string) error
}

// FileStatInfo 文件信息
type FileStatInfo struct {
	Name        string    // 文件名
	Size        int64     // 文件大小
	Ext         string    // 文件后缀
	MimeType    string    // 媒体类型
	LastModTime time.Time // 最后修改时间
}

// ObjectPart container for particular part of an object.
type ObjectPart struct {
	PartNumber     int
	ETag           string
	PartObjectName string
}

// LocalSystemConfig 本地存储 配置信息
type LocalSystemConfig struct {
	Root          string `json:"root" yaml:"root"`
	SSL           bool   `json:"ssl" yaml:"ssl"`
	BucketPublic  string `json:"bucket_public" yaml:"bucket_public"`
	BucketPrivate string `json:"bucket_private" yaml:"bucket_private"`
	Endpoint      string `json:"endpoint" yaml:"endpoint"`
}

// MinioSystemConfig 私有化 Minio 配置信息
type MinioSystemConfig struct {
	SSL           bool   `json:"ssl" yaml:"ssl"`
	SecretId      string `json:"secret_id" yaml:"secret_id"`
	SecretKey     string `json:"secret_key" yaml:"secret_key"`
	BucketPublic  string `json:"bucket_public" yaml:"bucket_public"`
	BucketPrivate string `json:"bucket_private" yaml:"bucket_private"`
	Endpoint      string `json:"endpoint" yaml:"endpoint"`
}
