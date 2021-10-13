package helper

import (
	"crypto/md5"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"time"
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

// GenImageName 随机生成指定后缀的图片名
func GenImageName(ext string, width, height int) string {
	str := fmt.Sprintf("%d%s", time.Now().Unix(), GetRandomString(10))

	return fmt.Sprintf("%x_%dx%d.%s", md5.Sum([]byte(str)), width, height, ext)
}
