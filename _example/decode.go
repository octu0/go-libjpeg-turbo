package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"

	"github.com/octu0/go-libjpeg-turbo"
)

func main() {
	f, err := os.Open("./testdata/src.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	decoder, err := turbojpeg.CreateDecoder()
	if err != nil {
		panic(err)
	}
	defer decoder.Close()

	header, err := decoder.DecodeHeader(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", header)

	img, err := decoder.Decode(data, turbojpeg.PixelFormatRGBA)
	if err != nil {
		panic(err)
	}
	defer img.Close()

	rgba := &image.RGBA{
		Pix:    img.Bytes(),
		Stride: img.Width * 4,
		Rect:   image.Rect(0, 0, img.Width, img.Height),
	}

	savePng(rgba)

	ref, err := decoder.DecodeToRGBA(data)
	if err != nil {
		panic(err)
	}
	defer ref.Close()

	savePng(ref.Image)
}

func savePng(rgba *image.RGBA) {
	out, err := os.CreateTemp("/tmp", "out*.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if err := png.Encode(out, rgba); err != nil {
		panic(err)
	}
	fmt.Println("png out =", out.Name())
}
