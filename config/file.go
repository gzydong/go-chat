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

// QiniuSystem 七牛云文件存储
type QiniuSystem struct {
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	Domain    string `json:"domain" yaml:"domain"`
}

// OssSystem 阿里云 OSS 文件存储
type CosSystem struct {
	SecretId  string `json:"secret_id" yaml:"secret_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	Region    string `json:"region" yaml:"region"`
}

type Filesystem struct {
	Default string      `json:"default" yaml:"default"`
	Local   LocalSystem `json:"local" yaml:"local"`
	Oss     OssSystem   `json:"oss" yaml:"oss"`
	Qiniu   QiniuSystem `json:"qiniu" yaml:"qiniu"`
	Cos     CosSystem   `json:"cos" yaml:"cos"`
}
