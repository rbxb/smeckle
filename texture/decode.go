// Package texture is used to unpack WebP textures.
package texture

import (
	"image"
	"image/color"
	"io"

	"golang.org/x/image/webp"
)

func convertNYCbCrA(c color.NYCbCrA) color.RGBA {
	r, g, b := color.YCbCrToRGB(c.Y, c.Cb, c.Cr)
	return color.RGBA{r, g, b, 255}
}

func convertNRGBA(c color.NRGBA) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, 255}
}

// Decode reads a WebP from the reader and returns color and specular image objects.
func Decode(r io.ReadSeeker) (image.Image, image.Image) {
	/*
		// Skip 48 bytes.
		// I don't know what these bytes are for.
		r.Seek(48,0)

		// Read the number of mipmaps in the file.
		var mipCount uint32
		if err := binary.Read(r, binary.LittleEndian, &mipCount); err != nil {
			panic(err)
		}
		if mipCount > 20 {
			mipCount = 20
		}

		// Skip 28 bytes.
		// I don't know what these bytes are for.
		r.Seek(28,0)

		// Read the size of the mipmaps.
		var mipSize uint32
		if err := binary.Read(r, binary.LittleEndian, &mipSize); err != nil {
			panic(err)
		}
	*/

	// I don't care about mipmaps. I only need the first one.
	r.Seek(84, 0)

	// Read the WebP data.
	img, err := webp.Decode(r)
	if err != nil {
		panic(err)
	}

	// Split the texture into a color texture and a specular map.
	colorImg := image.NewRGBA(img.Bounds())
	specImg := image.NewGray(img.Bounds())
	// Iterate through each pixel in the texture.
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			pixel := img.At(x, y)
			var colorPixel color.RGBA
			var specPixel color.Gray
			switch img.ColorModel() {
			// Dumb color conversion stuff.
			case color.NYCbCrAModel:
				colorPixel = convertNYCbCrA(pixel.(color.NYCbCrA))
				specPixel = color.Gray{pixel.(color.NYCbCrA).A}
			case color.NRGBAModel:
				colorPixel = convertNRGBA(pixel.(color.NRGBA))
				specPixel = color.Gray{pixel.(color.NRGBA).A}
			default:
				colorPixel = color.RGBAModel.Convert(pixel).(color.RGBA)
				specPixel = color.Gray{255}
			}
			colorImg.Set(x, y, colorPixel)
			specImg.Set(x, y, specPixel)
		}
	}
	return colorImg, specImg
}
