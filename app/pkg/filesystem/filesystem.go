package filesystem

import (
	"go-chat/config"
	"time"
)

type AdapterInterface interface {
	// Write 文件写入
	Write(data []byte, filePath string) error

	// WriteLocal 本地文件上传
	WriteLocal(localFile string, filePath string) error

	Copy(srcPath, filePath string) error

	// Delete 删除一个文件或空文件夹
	Delete(filePath string) error
	DeleteDir(path string) error
	CreateDir(path string) error

	// Stat 文件信息
	Stat(filePath string) (*FileStat, error)
}

// FileStat 文件信息
type FileStat struct {
	Name        string
	Size        int64
	Ext         string
	LastModTime time.Time
	MimeType    string
}

type Filesystem struct {
	AdapterInterface
}

func NewFilesystem(conf *config.Config) *Filesystem {
	var driver AdapterInterface
	switch conf.Filesystem.Default {
	case "oss":
		driver = NewOssFilesystem(conf)
		break
	case "qiniu":
		driver = NewQiniuFilesystem(conf)
		break
	default:
		driver = NewLocalFilesystem(conf)
	}

	return &Filesystem{driver}
}
