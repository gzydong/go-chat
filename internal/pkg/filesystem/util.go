package filesystem

import (
	"io/ioutil"
	"mime/multipart"
)

func ReadMultipartStream(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer src.Close()

	return ioutil.ReadAll(src)
}
