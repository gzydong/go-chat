package provider

import (
	"go-chat/config"
	"go-chat/internal/pkg/filesystem"
)

func NewFilesystem(conf *config.Config) *filesystem.Filesystem {

	s := &filesystem.Filesystem{}

	s.SetDriver(conf.Filesystem.Default)

	s.Local = filesystem.NewLocalFilesystem(conf)
	s.Cos = filesystem.NewCosFilesystem(conf)

	switch s.Driver() {
	case "cos":
		s.Default = s.Cos
	default:
		s.Default = s.Local
	}

	return s
}
