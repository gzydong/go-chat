package utils

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type MetaImage struct {
	Width  int
	Height int
}

func ReadFileImage(r io.Reader) *MetaImage {
	c, _, _ := image.DecodeConfig(r)

	return &MetaImage{c.Width, c.Height}
}
