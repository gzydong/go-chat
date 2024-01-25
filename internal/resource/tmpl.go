package resource

import (
	"embed"
)

//go:embed "templates"
var templates embed.FS

func Template() embed.FS {
	return templates
}
