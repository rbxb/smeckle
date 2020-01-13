// Package png is used to write an image as a png.
package png

import (
	"image"
	"image/png"
	"io"
)

// Encode writes an image as a png.
func Encode(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}
