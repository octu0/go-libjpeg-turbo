package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	"github.com/octu0/go-libjpeg-turbo"
)

func main() {
	f, err := os.Open("./testdata/src.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rgba, err := PNGToRGBA(f)
	if err != nil {
		panic(err)
	}

	out, err := os.CreateTemp("/tmp", "out*.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	encoder, err := turbojpeg.CreateEncoder()
	if err != nil {
		panic(err)
	}
	defer encoder.Close()

	if _, err := encoder.EncodeRGBA(out, rgba, 85); err != nil {
		panic(err)
	}

	fmt.Println("jpg out =", out.Name())
}

func PNGToRGBA(r io.Reader) (*image.RGBA, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	if i, ok := img.(*image.RGBA); ok {
		return i, nil
	}

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y += 1 {
		for x := b.Min.X; x < b.Max.X; x += 1 {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			rgba.Set(x, y, c)
		}
	}
	return rgba, nil
}
