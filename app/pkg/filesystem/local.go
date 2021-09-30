package filesystem

import (
	"fmt"
	"io"
	"os"
	"path"
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

func (s *LocalFilesystem) path(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.conf.Filesystem.Local.Root, "/"), strings.TrimLeft(path, "/"))
}

// 判断目录是否存在
func isDirExist(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}

	return s.IsDir()
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

func (s *LocalFilesystem) Copy(srcPath, filePath string) error {
	return s.WriteLocal(s.path(srcPath), filePath)
}

func (s *LocalFilesystem) Delete(filePath string) error {
	return os.Remove(s.path(filePath))
}

func (s *LocalFilesystem) CreateDir(dir string) error {
	return os.MkdirAll(s.path(dir), 055)
}

func (s *LocalFilesystem) DeleteDir(dir string) error {
	return os.RemoveAll(s.path(dir))
}

func (s *LocalFilesystem) Stat(filePath string) {
	info, _ := os.Stat(s.path(filePath))

	fmt.Printf("%#v", info)
}
