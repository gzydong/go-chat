package config

import "go-chat/internal/pkg/filesystem"

type Filesystem struct {
	Default string                       `json:"default" yaml:"default"`
	Local   filesystem.LocalSystemConfig `json:"local" yaml:"local"`
	Minio   filesystem.MinioSystemConfig `json:"minio" yaml:"minio"`
}
