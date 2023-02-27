package turbojpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include -I/usr/local/include -I/usr/include -I/opt/libjpeg-turbo/include
#cgo LDFLAGS: -L${SRCDIR}/lib -L/usr/local/lib -L/usr/lib -L/opt/libjpeg-turbo/lib64 -L/opt/libjpeg-turbo/lib -lturbojpeg -lm -ldl
#include <stdint.h>
#include <stdlib.h>

#include "decode.h"
*/
import "C"

import (
	"bytes"
	"image"
	"runtime"
	"runtime/cgo"
	"sync/atomic"
	"unsafe"

	"github.com/octu0/cgobytepool"
	"github.com/pkg/errors"
)

type Header struct {
	Width       int
	Height      int
	Subsampling Subsampling
	ColorSpace  ColorSpace
}

type Image struct {
	Header
	Format    PixelFormat
	result    *C.jpeg_decode_result_t
	closed    int32
	closeFunc func()
}

func (i *Image) Bytes() []byte {
	return cgobytepool.GoBytes(unsafe.Pointer(i.result.data), int(i.result.data_size))
}

func (i *Image) Close() {
	if atomic.CompareAndSwapInt32(&i.closed, 0, 1) {
		runtime.SetFinalizer(i, nil)

		i.closeFunc()
	}
}

func finalizeImage(i *Image) {
	i.Close()
}

type ImageRGBARef struct {
	Header
	Format    PixelFormat
	Image     *image.RGBA
	buf       *bytes.Buffer
	closed    int32
	closeFunc func()
}

func (i *ImageRGBARef) Close() {
	if atomic.CompareAndSwapInt32(&i.closed, 0, 1) {
		runtime.SetFinalizer(i, nil)

		i.closeFunc()
	}
}

func finalizeImageRGBARef(i *ImageRGBARef) {
	i.Close()
}

type Decoder struct {
	pool   cgobytepool.Pool
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

func (d *Decoder) decode(data []byte, format PixelFormat) (*Image, error) {
	ctx := cgobytepool.CgoHandle(d.pool)

	r := unsafe.Pointer(C.decode_jpeg(
		unsafe.Pointer(&ctx),
		(C.tjhandle)(d.handle),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.ulong(len(data)),
		C.int(format),
	))
	if r == nil {
		ctx.Delete()
		return nil, errors.Errorf("failed to call tjDecompress2()")
	}

	result := (*C.jpeg_decode_result_t)(r)
	header := Header{
		Width:       int(result.width),
		Height:      int(result.height),
		Subsampling: Subsampling(result.subsampling),
		ColorSpace:  ColorSpace(result.colorspace),
	}
	img := &Image{
		Header:    header,
		Format:    format,
		result:    result,
		closed:    0,
		closeFunc: createImageCloseFunc(ctx, result),
	}
	return img, nil
}

func (d *Decoder) Decode(data []byte, format PixelFormat) (*Image, error) {
	img, err := d.decode(data, format)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	runtime.SetFinalizer(img, finalizeImage)
	return img, nil
}

func (d *Decoder) DecodeToRGBA(data []byte) (*ImageRGBARef, error) {
	img, err := d.decode(data, PixelFormatRGBA)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rgba := &image.RGBA{
		Pix:    img.Bytes(),
		Stride: img.Header.Width * 4,
		Rect:   image.Rect(0, 0, img.Header.Width, img.Header.Height),
	}
	ref := &ImageRGBARef{
		Header:    img.Header,
		Format:    PixelFormatRGBA,
		Image:     rgba,
		closed:    0,
		closeFunc: createImageRGBARefCloseFunc(img),
	}
	runtime.SetFinalizer(ref, finalizeImageRGBARef)
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
	d, err := CreateDecoderWithPool(defaultCgoBytePool)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return d, nil
}

func CreateDecoderWithPool(p cgobytepool.Pool) (*Decoder, error) {
	h := unsafe.Pointer(C.tjInitDecompress())
	if h == nil {
		return nil, errors.Errorf("failed to call tjInitDecompress()")
	}
	d := &Decoder{
		pool:   p,
		handle: h,
		closed: 0,
	}
	runtime.SetFinalizer(d, finalizeDecoder)
	return d, nil
}

func createImageCloseFunc(ctx cgo.Handle, result *C.jpeg_decode_result_t) func() {
	return func() {
		defer ctx.Delete()

		C.free_jpeg_decode_result(unsafe.Pointer(&ctx), result)
	}
}

func createImageRGBARefCloseFunc(img *Image) func() {
	return func() {
		img.Close()
	}
}
