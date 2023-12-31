package strutil

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenImageName 随机生成指定后缀的图片名
func GenImageName(ext string, width, height int) string {
	return fmt.Sprintf("%s_%dx%d.%s", uuid.New().String(), width, height, ext)
}

func GenFileName(ext string) string {
	return fmt.Sprintf("%s.%s", uuid.New().String(), ext)
}

func GenMediaObjectName(ext string, width, height int) string {
	var (
		mediaType = "common"
		fileName  = GenFileName(ext)
	)

	switch ext {
	case "png", "jpeg", "jpg", "gif", "webp", "svg", "ico":
		mediaType = "image"
		fileName = GenImageName(ext, width, height)
	case "mp3", "wav", "aac", "ogg", "flac":
		mediaType = "audio"
	case "mp4", "avi", "mov", "wmv", "mkv":
		mediaType = "video"
	}

	return fmt.Sprintf("media/%s/%s/%s", mediaType, time.Now().Format("200601"), fileName)
}
