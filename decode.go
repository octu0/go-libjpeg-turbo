package turbojpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include -I/usr/local/include -I/usr/include -I/opt/libjpeg-turbo/include
#cgo LDFLAGS: -L${SRCDIR} -L/usr/local/lib -L/usr/lib -L/opt/libjpeg-turbo/lib64 -L/opt/libjpeg-turbo/lib -lturbojpeg -lm -ldl
#include <stdint.h>
#include <stdlib.h>

#include "decode.h"
*/
import "C"

import (
	"bytes"
	"image"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/pkg/errors"
)

var (
	imageBufPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 16*1024))
		},
	}
)

type Header struct {
	Width       int
	Height      int
	Subsampling Subsampling
	ColorSpace  ColorSpace
}

type Image struct {
	Header
	Format PixelFormat
	buf    *bytes.Buffer
	closed int32
}

func (i *Image) Bytes() []byte {
	return i.buf.Bytes()
}

func (i *Image) Close() {
	if atomic.CompareAndSwapInt32(&i.closed, 0, 1) {
		runtime.SetFinalizer(i, nil)

		imageBufPool.Put(i.buf)
	}
}

type ImageRGBARef struct {
	Header
	Format PixelFormat
	Image  *image.RGBA
	buf    *bytes.Buffer
	closed int32
}

func (i *ImageRGBARef) Close() {
	if atomic.CompareAndSwapInt32(&i.closed, 0, 1) {
		runtime.SetFinalizer(i, nil)

		imageBufPool.Put(i.buf)
	}
}

type Decoder struct {
	handle unsafe.Pointer // tjhandle
	closed int32
}

func (d *Decoder) DecodeHeader(data []byte) (Header, error) {
	h := unsafe.Pointer(C.decode_jpeg_header(
		(C.tjhandle)(d.handle),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.ulong(len(data)),
	))
	if h == nil {
		return Header{}, errors.Errorf("failed to call tjDecompressHeader3()")
	}
	header := (*C.jpeg_header_t)(h)
	defer C.free_jpeg_header(header)

	return Header{
		Width:       int(header.width),
		Height:      int(header.height),
		Subsampling: Subsampling(header.subsampling),
		ColorSpace:  ColorSpace(header.colorspace),
	}, nil
}

func (d *Decoder) decode(out *bytes.Buffer, data []byte, format PixelFormat) (Header, error) {
	r := unsafe.Pointer(C.decode_jpeg(
		(C.tjhandle)(d.handle),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.ulong(len(data)),
		C.int(format),
	))
	if r == nil {
		return Header{}, errors.Errorf("failed to call tjDecompress2()")
	}
	result := (*C.jpeg_decode_result_t)(r)
	defer C.free_jpeg_decode_result(result)

	out.Reset()
	out.Grow(int(result.data_size))
	out.Write(C.GoBytes(unsafe.Pointer(result.data), result.data_size)) // copy C allocated

	return Header{
		Width:       int(result.width),
		Height:      int(result.height),
		Subsampling: Subsampling(result.subsampling),
		ColorSpace:  ColorSpace(result.colorspace),
	}, nil
}

func (d *Decoder) Decode(data []byte, format PixelFormat) (*Image, error) {
	out := imageBufPool.Get().(*bytes.Buffer)
	header, err := d.decode(out, data, format)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	img := &Image{
		Header: header,
		Format: format,
		buf:    out,
		closed: 0,
	}
	runtime.SetFinalizer(img, func(me *Image) {
		me.Close()
	})
	return img, nil
}

func (d *Decoder) DecodeToRGBA(data []byte) (*ImageRGBARef, error) {
	out := imageBufPool.Get().(*bytes.Buffer)
	out.Reset()
	header, err := d.decode(out, data, PixelFormatRGBA)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rgba := &image.RGBA{
		Pix:    out.Bytes(),
		Stride: header.Width * 4,
		Rect:   image.Rect(0, 0, header.Width, header.Height),
	}
	ref := &ImageRGBARef{
		Header: header,
		Format: PixelFormatRGBA,
		Image:  rgba,
		buf:    out,
		closed: 0,
	}
	runtime.SetFinalizer(ref, func(me *ImageRGBARef) {
		me.Close()
	})
	return ref, nil
}

func (d *Decoder) Close() {
	if atomic.CompareAndSwapInt32(&d.closed, 0, 1) {
		runtime.SetFinalizer(d, nil)

		C.tjDestroy((C.tjhandle)(d.handle))
	}
}

func finalizeDecoder(d *Decoder) {
	d.Close()
}

func CreateDecoder() (*Decoder, error) {
	h := unsafe.Pointer(C.tjInitDecompress())
	if h == nil {
		return nil, errors.Errorf("failed to call tjInitDecompress()")
	}
	d := &Decoder{
		handle: h,
		closed: 0,
	}
	runtime.SetFinalizer(d, finalizeDecoder)
	return d, nil
}
