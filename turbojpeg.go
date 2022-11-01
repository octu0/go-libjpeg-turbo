package turbojpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include -I/usr/local/include -I/usr/include -I/opt/libjpeg-turbo/include
#cgo LDFLAGS: -L${SRCDIR} -L/usr/local/lib -L/usr/lib -L/opt/libjpeg-turbo/lib64 -L/opt/libjpeg-turbo/lib -lturbojpeg -lm -ldl

#include "turbojpeg.h"
*/
import "C"

import (
	"bytes"
	"sync"
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

var (
	imageBufPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 16*1024))
		},
	}
)
