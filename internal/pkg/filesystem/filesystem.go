package filesystem

import (
	"time"

	"go-chat/config"
)

type IAdapter interface {
	// Write 文件写入
	Write(data []byte, filePath string) error

	// WriteLocal 本地文件上传
	WriteLocal(localFile string, filePath string) error

	// Copy 文件拷贝
	Copy(srcPath, filePath string) error

	// Delete 删除一个文件或空文件夹
	Delete(filePath string) error

	// DeleteDir 删除文件夹
	DeleteDir(path string) error

	// CreateDir 递归创建文件夹
	CreateDir(path string) error

	// Stat 文件信息
	Stat(filePath string) (*FileStat, error)

	// PublicUrl 获取公开文件的访问地址
	PublicUrl(filePath string) string

	// PrivateUrl 获取私有文件的访问地址
	PrivateUrl(filePath string, timeout int) string

	// ReadStream 读取文件内容
	ReadStream(filePath string) ([]byte, error)

	InitiateMultipartUpload(filePath string, fileName string) (string, error)
}

// FileStat 文件信息
type FileStat struct {
	Name        string    // 文件名
	Size        int64     // 文件大小
	Ext         string    // 文件后缀
	LastModTime time.Time // 最后修改时间
	MimeType    string    // 媒体类型
}

type Filesystem struct {
	driver  string
	Default IAdapter
	Local   *LocalFilesystem
	Cos     *CosFilesystem
}

func NewFilesystem(conf *config.Config) *Filesystem {
	s := &Filesystem{}

	s.driver = conf.Filesystem.Default

	s.Local = NewLocalFilesystem(conf)
	s.Cos = NewCosFilesystem(conf)

	switch s.driver {
	case "cos":
		s.Default = s.Cos
	default:
		s.Default = s.Local
	}

	return s
}

func (f *Filesystem) Driver() string {
	return f.driver
}

func (f *Filesystem) SetDriver(value string) {
	f.driver = value
}
