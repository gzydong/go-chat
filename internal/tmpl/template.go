package tmpl

import (
	"embed"
)

//go:embed "resource"
var templates embed.FS

func Templates() embed.FS {
	return templates
}
