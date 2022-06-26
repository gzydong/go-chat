package resource

import (
	"embed"
)

//go:embed "templates"
var templates embed.FS

func Templates() embed.FS {
	return templates
}
