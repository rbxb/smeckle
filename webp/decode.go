package webp

import (
	"io"
	"encoding/binary"
	"image"
	"image/color"
	"golang.org/x/image/webp"
)

func convertNYCbCrA(c color.NYCbCrA) color.RGBA {
	r,g,b := color.YCbCrToRGB(c.Y, c.Cb, c.Cr)
	return color.RGBA{r, g, b, 255}
}

func convertNRGBA(c color.NRGBA) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, 255}
}

func Decode(r io.ReadSeeker) (image.Image, image.Image, error) {
	r.Seek(48,0)

	var mipCount uint32
	if err := binary.Read(r, binary.LittleEndian, &mipCount); err != nil {
		panic(err)
	}
	if mipCount > 20 {
		mipCount = 20
	}
	//Read the number of mips in the file.

	r.Seek(28,0)

	var mipSize uint32
	if err := binary.Read(r, binary.LittleEndian, &mipSize); err != nil {
		panic(err)
	}
	//Read the mip size.

	img, err := webp.Decode(r)
	if err != nil {
		panic(err)
	}
	//Read the image data.

	colorImg := image.NewRGBA(img.Bounds())
	specImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			pixel := img.At(x,y)
			var colorPixel color.RGBA
			var specPixel color.Gray
			switch img.ColorModel() {
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
			colorImg.Set(x,y,colorPixel)
			specImg.Set(x,y,specPixel)
		}
	}
	return colorImg, specImg, nil
}