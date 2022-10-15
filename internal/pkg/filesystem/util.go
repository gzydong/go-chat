package filesystem

import (
	"io"
	"mime/multipart"
)

func ReadMultipartStream(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer src.Close()

	return io.ReadAll(src)
}
