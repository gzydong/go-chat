package config

// LocalSystem 本地存储
type LocalSystem struct {
	Root   string `json:"root" yaml:"root"`
	Domain string `json:"domain" yaml:"domain"`
}

// OssSystem 阿里云 OSS 文件存储
type OssSystem struct {
	AccessID     string `json:"access_id" yaml:"access_id"`
	AccessSecret string `json:"access_secret" yaml:"access_secret"`
	Bucket       string `json:"bucket" yaml:"bucket"`
	Endpoint     string `json:"endpoint" yaml:"endpoint"`
}

type Filesystem struct {
	Driver string      `json:"driver" yaml:"driver"`
	Local  LocalSystem `json:"local" yaml:"local"`
	Oss    OssSystem   `json:"oss" yaml:"oss"`
}
