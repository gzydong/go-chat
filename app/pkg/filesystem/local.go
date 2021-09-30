package filesystem

import (
	"fmt"
	"go-chat/config"
	"os"
	"path"
	"strings"
)

type LocalFilesystem struct {
	conf *config.Config
}

func NewLocalFilesystem(conf *config.Config) *LocalFilesystem {
	return &LocalFilesystem{
		conf: conf,
	}
}

func (s *LocalFilesystem) Write(data []byte, filePath string) {
	filePath = s.path(filePath)

	_ = os.MkdirAll(path.Dir(filePath), 0777)

	f, _ := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)

	_, _ = f.Write(data)
}

func (s *LocalFilesystem) WriteLocal(localFile string, filePath string) {

}

func (s *LocalFilesystem) Update() {
	fmt.Println("LocalFilesystem :", "Update")
}

func (s *LocalFilesystem) Rename() {

}

func (s *LocalFilesystem) Copy() {

}

func (s *LocalFilesystem) Delete(filePath string) error {
	return os.Remove(s.path(filePath))
}

func (s *LocalFilesystem) DeleteDir(dir string) error {
	return os.RemoveAll(s.path(dir))
}

func (s *LocalFilesystem) CreateDir(dir string) error {
	return os.MkdirAll(s.path(dir), 0777)
}

func (s *LocalFilesystem) UploadFile() {

}

func (s *LocalFilesystem) path(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.conf.Filesystem.Local.Root, "/"), strings.TrimLeft(path, "/"))
}
