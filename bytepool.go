package turbojpeg

/*
#include "bytepool.h"
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
