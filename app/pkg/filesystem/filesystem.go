package filesystem

import (
	"go-chat/config"
)

type AdapterInterface interface {
	// Write 文件写入
	Write(data []byte, filePath string)
	// WriteLocal 本地文件上传
	WriteLocal(localFile string, filePath string)

	Update()
	Rename()
	Copy()

	// Delete 删除一个文件或空文件夹
	Delete(filePath string) error
	DeleteDir(path string) error
	CreateDir(path string) error
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
