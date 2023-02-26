package turbojpeg

import (
	"bytes"
	"image/jpeg"
	"io"
	"os"
	"testing"
)

func BenchmarkJpegDecode(b *testing.B) {
	jpegfile, err := os.Open("./testdata/src.jpg")
	if err != nil {
		b.Fatalf("no error: %+v", err)
	}
	defer jpegfile.Close()

	data, err := io.ReadAll(jpegfile)
	if err != nil {
		b.Fatalf("no error: %+v", err)
	}

	b.Run("image/jpeg.Decode", func(tb *testing.B) {
		for i := 0; i < tb.N; i += 1 {
			if _, err := jpeg.Decode(bytes.NewReader(data)); err != nil {
				tb.Fatalf("no error: %+v", err)
			}
		}
	})
	b.Run("turbojpeg.Decode", func(tb *testing.B) {
		decoder, err := CreateDecoder()
		if err != nil {
			tb.Fatalf("no error: %+v", err)
		}

		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			img, err := decoder.Decode(data, PixelFormatRGBA)
			if err != nil {
				tb.Fatalf("no error: %+v", err)
			}
			img.Close()
		}
	})
	b.Run("turbojpeg.DecodeToRGBA", func(tb *testing.B) {
		decoder, err := CreateDecoder()
		if err != nil {
			tb.Fatalf("no error: %+v", err)
		}

		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			ref, err := decoder.DecodeToRGBA(data)
			if err != nil {
				tb.Fatalf("no error: %+v", err)
			}
			ref.Close()
		}
	})
}
