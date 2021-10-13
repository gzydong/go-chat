package filesystem

import (
	"time"

	"go-chat/config"
)

type AdapterInterface interface {
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

	PublicUrl(filePath string) string

	PrivateUrl(filePath string, timeout int) string
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
