package helper

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

// ReadImage 读取文件大小
func ReadImage(src string) map[string]int {
	file, _ := os.Open(src)

	defer file.Close()

	c, _, _ := image.DecodeConfig(file)

	return map[string]int{
		"width":  c.Width,
		"height": c.Height,
	}
}

func ReadFileImage(r io.Reader) map[string]int {
	c, _, _ := image.DecodeConfig(r)

	return map[string]int{
		"width":  c.Width,
		"height": c.Height,
	}
}
