package provider

import (
	"github.com/gzydong/go-chat/config"
	"github.com/gzydong/go-chat/internal/pkg/filesystem"
)

func NewFilesystem(conf *config.Config) filesystem.IFilesystem {
	if conf.Filesystem.Default == filesystem.MinioDriver {
		return filesystem.NewMinioFilesystem(conf.Filesystem.Minio)
	}

	return filesystem.NewLocalFilesystem(conf.Filesystem.Local)
}
