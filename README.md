# `go-libjpeg-turbo`

[![License](https://img.shields.io/github/license/octu0/go-libjpeg-turbo)](https://github.com/octu0/go-libjpeg-turbo/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/octu0/go-libjpeg-turbo?status.svg)](https://godoc.org/github.com/octu0/go-libjpeg-turbo)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/go-libjpeg-turbo)](https://goreportcard.com/report/github.com/octu0/go-libjpeg-turbo)
[![Releases](https://img.shields.io/github/v/release/octu0/go-libjpeg-turbo)](https://github.com/octu0/go-libjpeg-turbo/releases)

Go bindings for [libjpeg-turbo](https://github.com/libjpeg-turbo/libjpeg-turbo)  

## Requirements

requires libjpeg-turbo [install](https://github.com/libjpeg-turbo/libjpeg-turbo) on your system

## Usage

### Decode

```go
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
	data, err := readFile(“/path/to/image.jpg”)
	if err != nil {
		panic(err)
	}

	// create Decoder
	decoder, err := turbojpeg.CreateDecoder()
	if err != nil {
		panic(err)
	}
	defer decoder.Close()

	// decode Header only
	header, err := decoder.DecodeHeader(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", header)
	// => {Width:320 Height:240 Subsampling:4:2:0 ColorSpace:YCbCr}

	// decode JPEG image To RGBA
	imgRef, err := decoder.DecodeToRGBA(data)
	if err != nil {
		panic(err)
	}
	defer imgRef.Close()

	saveImage(imgRef.Image)

	// or Decode with PixelFormat
	img, err := decoder.Decode(data, turbojpeg.PixelFormatRGBA)
	if err != nil {
		panic(err)
	}
	defer img.Close()

	rgba := &image.RGBA{
		Pix:    img.Bytes(),   // decoded image
		Stride: img.Width * 4, // 4 = r + g + b + a
		Rect:   image.Rect(0, 0, img.Width, img.Height),
	}

	saveImage(rgba)
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open("./testdata/src.jpg")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func saveImage(img *image.RGBA)  {
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
```

### Encode

```go
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
	rgba, err := readFileRGBA("./testdata/src.png")
	if err != nil {
		panic(err)
	}

	out, err := os.CreateTemp("/tmp", "out*.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// create Encoder
	encoder, err := turbojpeg.CreateEncoder()
	if err != nil {
		panic(err)
	}
	defer encoder.Close()

	// encode jpeg From *image.RGBA
	if _, err := encoder.EncodeRGBA(out, rgba, 85); err != nil {
		panic(err)
	}

	fmt.Println("jpg out =", out.Name())
}

func readFileRGBA(path string) (*image.RGBA, error) {
	f, err := os.Open("./testdata/src.png")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rgba, err := PNGToRGBA(f)
	if err != nil {
		return nil, err
	}
	return rgba, nil
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
```

## Benchmark

### Decode

~3x faster than [image/jpeg.Decode](https://pkg.go.dev/image/jpeg#Decode)

```
goos: darwin
goarch: amd64
pkg: github.com/octu0/go-libjpeg-turbo
cpu: Intel(R) Core(TM) i5-8210Y CPU @ 1.60GHz
BenchmarkJpegDecode
BenchmarkJpegDecode/image/jpeg.Decode
BenchmarkJpegDecode/image/jpeg.Decode-4         	     429	   2568782 ns/op	  136976 B/op	      19 allocs/op
BenchmarkJpegDecode/turbojpeg.Decode
BenchmarkJpegDecode/turbojpeg.Decode-4          	    1468	    806084 ns/op	     775 B/op	      19 allocs/op
BenchmarkJpegDecode/turbojpeg.DecodeToRGBA
BenchmarkJpegDecode/turbojpeg.DecodeToRGBA-4    	    1460	    810547 ns/op	     920 B/op	      23 allocs/op
PASS
```

### Encode

~9x faster than [image/jpeg.Encode](https://pkg.go.dev/image/jpeg#Encode)

```
goos: darwin
goarch: amd64
pkg: github.com/octu0/go-libjpeg-turbo
cpu: Intel(R) Core(TM) i5-8210Y CPU @ 1.60GHz
BenchmarkJpegEncode
BenchmarkJpegEncode/image/jpeg.Encode
BenchmarkJpegEncode/image/jpeg.Encode-4         	     422	   3012564 ns/op	    4400 B/op	       4 allocs/op
BenchmarkJpegEncode/turbojpeg.Encode
BenchmarkJpegEncode/turbojpeg.Encode-4          	    3394	    345214 ns/op	     343 B/op	       8 allocs/op
PASS
```
