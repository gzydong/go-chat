package filesystem

import (
	"fmt"
	"go-chat/config"
)

type OssFilesystem struct {
}

func NewOssFilesystem(conf *config.Config) *OssFilesystem {
	return &OssFilesystem{}
}

func (s *OssFilesystem) Write(data []byte, filePath string) {

}

func (s *OssFilesystem) WriteLocal(localFile string, filePath string) {

}

func (s *OssFilesystem) Update() {
	fmt.Println("OssFilesystem :", "Update")
}

func (s *OssFilesystem) Rename() {

}

func (s *OssFilesystem) Copy(srcPath, filePath string) {

}

func (s *OssFilesystem) Delete(filePath string) error {
	return nil
}

func (s *OssFilesystem) DeleteDir(path string) error {
	return nil
}

func (s *OssFilesystem) CreateDir(path string) error {
	return nil
}
