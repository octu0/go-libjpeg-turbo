package turbojpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include -I/usr/local/include -I/usr/include -I/opt/libjpeg-turbo/include
#cgo LDFLAGS: -L${SRCDIR}/lib -L/usr/local/lib -L/usr/lib -L/opt/libjpeg-turbo/lib64 -L/opt/libjpeg-turbo/lib -lturbojpeg -lm -ldl
#include <stdint.h>
#include <stdlib.h>

#include "encode.h"
*/
import "C"

import (
	"image"
	"io"
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/pkg/errors"
)

type Encoder struct {
	handle unsafe.Pointer // tjhandle
	closed int32
}

func (e *Encoder) EncodeRGBA(out io.Writer, img *image.RGBA, quality int) (int, error) {
	width, height := img.Rect.Dx(), img.Rect.Dy()
	r := unsafe.Pointer(C.encode_jpeg(
		(C.tjhandle)(e.handle),
		(*C.uchar)(unsafe.Pointer(&img.Pix[0])),
		C.int(width),
		C.int(height),
		C.int(img.Stride),
		C.int(PixelFormatRGBA),
		C.int(Subsampling420),
		C.int(quality),
	))
	if r == nil {
		return 0, errors.Errorf("failed to call tjCompress2()")
	}
	result := (*C.jpeg_encode_result_t)(r)
	defer C.free_jpeg_encode_result(result)

	return out.Write(C.GoBytes(unsafe.Pointer(result.data), result.data_size))
}

func (e *Encoder) Close() {
	if atomic.CompareAndSwapInt32(&e.closed, 0, 1) {
		runtime.SetFinalizer(e, nil)

		C.tjDestroy((C.tjhandle)(e.handle))
	}
}

func finalizeEncoder(e *Encoder) {
	e.Close()
}

func CreateEncoder() (*Encoder, error) {
	h := unsafe.Pointer(C.tjInitCompress())
	if h == nil {
		return nil, errors.Errorf("failed to call tjInitCompress()")
	}

	e := &Encoder{
		handle: h,
		closed: 0,
	}
	runtime.SetFinalizer(e, finalizeEncoder)
	return e, nil
}
