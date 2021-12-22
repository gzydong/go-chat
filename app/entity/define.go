package entity

// 文件系统相关
const (
	FileDriveLocal = 1
	FileDriveCos   = 2
)

var fileDrives = map[string]int{
	"local": FileDriveLocal,
	"cos":   FileDriveCos,
}

func FileDriveMode(drive string) int {
	if val, ok := fileDrives[drive]; ok {
		return val
	}

	return 0
}
