package entity

const (
	MediaFileImage = 1 // 图片文件
	MediaFileVideo = 2 // 视频文件
	MediaFileAudio = 3 // 音频文件
	MediaFileOther = 4 // 其它文件
)

var mediaMaps = map[string]int{
	"gif":  MediaFileImage,
	"jpg":  MediaFileImage,
	"jpeg": MediaFileImage,
	"png":  MediaFileImage,
	"webp": MediaFileImage,
	"ogg":  MediaFileVideo,
	"mp3":  MediaFileVideo,
	"wav":  MediaFileVideo,
	"mp4":  MediaFileAudio,
	"webm": MediaFileAudio,
}

func GetMediaType(ext string) int {
	if val, ok := mediaMaps[ext]; ok {
		return val
	}

	return MediaFileOther
}
