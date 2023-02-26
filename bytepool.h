#include <stdlib.h>

extern void *turbojpeg_bytepool_get(void *context, size_t size);
extern void turbojpeg_bytepool_put(void *context, void *data, size_t size);
extern void turbojpeg_bytepool_free(void *context);
