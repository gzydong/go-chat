package entity

const (
	MediaFileImage = 1 // 图片文件
	MediaFileVideo = 2 // 视频文件
	MediaFileAudio = 3 // 音频文件
	MediaFileOther = 4 // 其它文件
)

var fileMediaMap = map[string]int{
	"gif":  MediaFileImage,
	"jpg":  MediaFileImage,
	"jpeg": MediaFileImage,
	"png":  MediaFileImage,
	"webp": MediaFileImage,
	"mp3":  MediaFileAudio,
	"wav":  MediaFileAudio,
	"mp4":  MediaFileVideo,
}

func GetMediaType(ext string) int {
	if val, ok := fileMediaMap[ext]; ok {
		return val
	}

	return MediaFileOther
}

// 文件系统相关
const (
	FileDriveLocal = 1
	FileDriveMinio = 2
)

var fileSystemDriveMap = map[string]int{
	"local": FileDriveLocal,
	"minio": FileDriveMinio,
}

func FileDriveMode(drive string) int {
	if val, ok := fileSystemDriveMap[drive]; ok {
		return val
	}

	return 0
}
