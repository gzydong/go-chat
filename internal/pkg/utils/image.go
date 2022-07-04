package utils

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type ImageMeta struct {
	Width  int
	Height int
}

func LoadImage(r io.Reader) *ImageMeta {
	c, _, _ := image.DecodeConfig(r)

	return &ImageMeta{c.Width, c.Height}
}
