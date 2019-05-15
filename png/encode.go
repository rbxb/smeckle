package png

import (
	"io"
	"image"
	"image/png"
)

func Encode(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}