package helper

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func ReadFileImage(r io.Reader) map[string]int {
	c, _, _ := image.DecodeConfig(r)

	return map[string]int{
		"width":  c.Width,
		"height": c.Height,
	}
}
