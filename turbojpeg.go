package turbojpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include -I/usr/local/include -I/usr/include -I/opt/libjpeg-turbo/include
#cgo LDFLAGS: -L${SRCDIR} -L/usr/local/lib -L/usr/lib -L/opt/libjpeg-turbo/lib64 -L/opt/libjpeg-turbo/lib -lturbojpeg -lm -ldl

#include "turbojpeg.h"
#include "pool.h"
*/
import "C"

import (
	"unsafe"

	"github.com/octu0/cgobytepool"
)

//export turbojpeg_bytepool_get
func turbojpeg_bytepool_get(ctx unsafe.Pointer, size C.size_t) unsafe.Pointer {
	return cgobytepool.HandlePoolGet(ctx, int(size))
}

//export turbojpeg_bytepool_put
func turbojpeg_bytepool_put(ctx unsafe.Pointer, data unsafe.Pointer, size C.size_t) {
	cgobytepool.HandlePoolPut(ctx, data, int(size))
}

//export turbojpeg_bytepool_free
func turbojpeg_bytepool_free(ctx unsafe.Pointer) {
	cgobytepool.HandlePoolFree(ctx)
}

var (
	defaultCgoBytePool = cgobytepool.NewPool(
		cgobytepool.DefaultMemoryAlignmentFunc,
		cgobytepool.WithPoolSize(1000, 4*1024),
		cgobytepool.WithPoolSize(1000, 8*1024),
		cgobytepool.WithPoolSize(1000, 16*1024),
		cgobytepool.WithPoolSize(1000, 32*1024),
		cgobytepool.WithPoolSize(1000, 64*1024),
	)
)

type Subsampling C.int

const (
	Subsampling444  Subsampling = C.TJSAMP_444
	Subsampling422  Subsampling = C.TJSAMP_422
	Subsampling420  Subsampling = C.TJSAMP_420
	SubsamplingGray Subsampling = C.TJSAMP_GRAY
	Subsampling440  Subsampling = C.TJSAMP_440
	Subsampling411  Subsampling = C.TJSAMP_411
)

type PixelFormat C.int

const (
	PixelFormatRGB  PixelFormat = C.TJPF_RGB
	PixelFormatBGR  PixelFormat = C.TJPF_BGR
	PixelFormatRGBX PixelFormat = C.TJPF_RGBX
	PixelFormatBGRX PixelFormat = C.TJPF_BGRX
	PixelFormatXBGR PixelFormat = C.TJPF_XBGR
	PixelFormatXRGB PixelFormat = C.TJPF_XRGB
	PixelFormatGray PixelFormat = C.TJPF_GRAY
	PixelFormatRGBA PixelFormat = C.TJPF_RGBA
	PixelFormatBGRA PixelFormat = C.TJPF_BGRA
	PixelFormatABGR PixelFormat = C.TJPF_ABGR
	PixelFormatARGB PixelFormat = C.TJPF_ARGB
	PixelFormatCMYK PixelFormat = C.TJPF_CMYK
)

type ColorSpace C.int

const (
	ColorSpaceRGB   ColorSpace = C.TJCS_RGB
	ColorSpaceYCbCr ColorSpace = C.TJCS_YCbCr
	ColorSpaceGray  ColorSpace = C.TJCS_GRAY
	ColorSpaceCMYK  ColorSpace = C.TJCS_CMYK
	ColorSpaceYCCK  ColorSpace = C.TJCS_YCCK
)
