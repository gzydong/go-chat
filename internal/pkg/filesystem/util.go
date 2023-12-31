package filesystem

import (
	"io"
	"mime/multipart"
	"os"
)

func ReadMultipartStream(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer src.Close()

	return io.ReadAll(src)
}

// isDirExist 判断目录是否存在
func isDirExist(fileAddr string) bool {
	s, err := os.Stat(fileAddr)

	return err == nil && s.IsDir()
}
