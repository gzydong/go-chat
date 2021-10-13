package filesystem

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go-chat/config"
)

type LocalFilesystem struct {
	conf *config.Config
}

func NewLocalFilesystem(conf *config.Config) *LocalFilesystem {
	return &LocalFilesystem{
		conf: conf,
	}
}

// path 获取文件地址绝对路径
func (s *LocalFilesystem) path(path string) string {
	return fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(s.conf.Filesystem.Local.Root, "/"),
		strings.TrimLeft(path, "/"),
	)
}

// isDirExist 判断目录是否存在
func isDirExist(fileAddr string) bool {
	s, err := os.Stat(fileAddr)

	return err == nil && s.IsDir()
}

// Write 上传 Byte 数组
func (s *LocalFilesystem) Write(data []byte, filePath string) error {
	filePath = s.path(filePath)

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

	_, err = f.Write(data)
	return err
}

// WriteLocal 本地文件上传
func (s *LocalFilesystem) WriteLocal(localFile string, filePath string) error {
	srcFile, err := os.Open(localFile)

	if err != nil {
		return err
	}

	defer srcFile.Close()

	dir := path.Dir(filePath)

	if len(dir) > 0 && !isDirExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	dstFile, err := os.OpenFile(s.path(filePath), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)

	return err
}

// Copy 文件拷贝
func (s *LocalFilesystem) Copy(srcPath, filePath string) error {
	return s.WriteLocal(s.path(srcPath), filePath)
}

// Delete 文件删除
func (s *LocalFilesystem) Delete(filePath string) error {
	return os.Remove(s.path(filePath))
}

// CreateDir 递归创建文件夹
func (s *LocalFilesystem) CreateDir(dir string) error {
	return os.MkdirAll(s.path(dir), 0755)
}

func (s *LocalFilesystem) DeleteDir(dir string) error {
	return os.RemoveAll(s.path(dir))
}

// Stat 文件信息
func (s *LocalFilesystem) Stat(filePath string) (*FileStat, error) {
	info, err := os.Stat(s.path(filePath))

	if err != nil {
		return nil, err
	}

	return &FileStat{
		Name:        filepath.Base(filePath),
		Size:        info.Size(),
		Ext:         filepath.Ext(filePath),
		MimeType:    "",
		LastModTime: info.ModTime(),
	}, nil
}

func (s *LocalFilesystem) PublicUrl(filePath string) string {
	return fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(s.conf.Filesystem.Local.Domain, "/"),
		strings.Trim(filePath, "/"),
	)
}

func (s *LocalFilesystem) PrivateUrl(filePath string, timeout int) string {
	return ""
}
