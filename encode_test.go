package turbojpeg

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"testing"
)

func BenchmarkJpegEncode(b *testing.B) {
	pngfile, err := os.Open("./testdata/src.png")
	if err != nil {
		b.Fatalf("no error: %+v", err)
	}
	defer pngfile.Close()

	rgba, err := testPNGToRGBA(pngfile)
	if err != nil {
		b.Fatalf("no error: %+v", err)
	}

	b.Run("image/jpeg.Encode", func(tb *testing.B) {
		for i := 0; i < tb.N; i += 1 {
			if err := jpeg.Encode(io.Discard, rgba, nil); err != nil {
				tb.Fatalf("no error: %+v", err)
			}
		}
	})
	b.Run("turbojpeg.Encode", func(tb *testing.B) {
		encoder, err := CreateEncoder()
		if err != nil {
			tb.Fatalf("no error: %+v", err)
		}
		defer encoder.Close()

		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			if _, err := encoder.EncodeRGBA(io.Discard, rgba, jpeg.DefaultQuality); err != nil {
				tb.Fatalf("no error: %+v", err)
			}
		}
	})
}

func testPNGToRGBA(r io.Reader) (*image.RGBA, error) {
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
