#include <stdio.h>
#include <string.h>
#include "turbojpeg.h"

#ifndef H_GO_LIBJPEG_TURBO_ENCODE
#define H_GO_LIBJPEG_TURBO_ENCODE

typedef struct jpeg_encode_result_t {
  unsigned char *data;
  int data_size;
} jpeg_encode_result_t;

void free_jpeg_encode_result(jpeg_encode_result_t *result) {
  if (NULL != result) {
    free(result->data);
  }
  free(result);
}

jpeg_encode_result_t *encode_jpeg(
  tjhandle handle,
  unsigned char *data,
  int width,
  int height,
  int stride,
  int src_pixel_format,
  int subsampling,
  int quality
) {
  unsigned char *out = NULL;
  unsigned long out_size = 0;
  int flags = 0;
  int ret = tjCompress2(
    handle,
    data,
    width,
    stride,
    height,
    src_pixel_format,
    &out,
    &out_size,
    subsampling,
    quality,
    flags
  );
  if (ret != 0) {
    return NULL;
  }

  jpeg_encode_result_t *result = (jpeg_encode_result_t*) malloc(sizeof(jpeg_encode_result_t));
  if(NULL == result) {
    tjFree(out);
    return NULL;
  }
  result->data = (unsigned char*) malloc(out_size);
  if(NULL == result->data) {
    free_jpeg_encode_result(result);
    tjFree(out);
    return NULL;
  }
  memcpy(result->data, out, out_size);
  tjFree(out);

  result->data_size = out_size;
  return result;
}

#endif
