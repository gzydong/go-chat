package testutil

import (
	"path/filepath"
	"runtime"

	"go-chat/config"
)

func GetConfig() *config.Config {
	_, file, _, _ := runtime.Caller(0)

	paths := []string{filepath.Dir(filepath.Dir(file)), "./config.yaml"}

	return config.New(filepath.Join(paths...))
}
