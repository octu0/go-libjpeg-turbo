#include <stdio.h>
#include <string.h>
#include "turbojpeg.h"

#ifndef H_GO_LIBJPEG_TURBO_DECODE
#define H_GO_LIBJPEG_TURBO_DECODE

typedef struct jpeg_header_t {
  int width;
  int height;
  int subsampling;
  int colorspace;
} jpeg_header_t;

typedef struct jpeg_decode_result_t {
  int width;
  int height;
  int subsampling;
  int colorspace;
  int pixel_format;
  unsigned char *data;
  int data_size;
} jpeg_decode_result_t;

void free_jpeg_header(jpeg_header_t *header) {
  free(header);
}

void free_jpeg_decode_result(jpeg_decode_result_t *result) {
  if(NULL != result) {
    free(result->data);
  }
  free(result);
}

jpeg_header_t *decode_jpeg_header(
  tjhandle handle,
  unsigned char *data,
  unsigned long data_size
) {
  jpeg_header_t *header = malloc(sizeof(jpeg_header_t));
  if(NULL == header) {
    return NULL;
  }
  memset(header, 0, sizeof(jpeg_header_t));

  int ret = tjDecompressHeader3(
    handle,
    data,
    data_size,
    &header->width,
    &header->height,
    &header->subsampling,
    &header->colorspace
  );
  if(0 != ret) {
    free_jpeg_header(header);
    return NULL;
  }
  return header;
}

jpeg_decode_result_t *decode_jpeg(
  tjhandle handle,
  unsigned char *data,
  unsigned long data_size,
  int dst_pixel_format
) {
  jpeg_decode_result_t *result = (jpeg_decode_result_t*) malloc(sizeof(jpeg_decode_result_t));
  if(NULL == result) {
    return NULL;
  }
  memset(result, 0, sizeof(jpeg_decode_result_t));

  int ret = tjDecompressHeader3(
    handle,
    data,
    data_size,
    &result->width,
    &result->height,
    &result->subsampling,
    &result->colorspace
  );
  if (0 != ret) {
    free_jpeg_decode_result(result);
    return NULL;
  }

  int pitch = tjPixelSize[dst_pixel_format] * result->width;
  int dst_size = pitch * result->height;
  result->data = (unsigned char*) malloc(dst_size);
  if (NULL == result->data) {
    free_jpeg_decode_result(result);
    return NULL;
  }
  result->data_size = dst_size;
  result->pixel_format = dst_pixel_format;

  int flags = 0;
  int ret_decompress = tjDecompress2(
    handle,
    data,
    data_size,
    result->data,
    result->width,
    pitch,
    result->height,
    dst_pixel_format,
    flags
  );
  if (0 != ret_decompress) {
    free_jpeg_decode_result(result);
    return NULL;
  }
  return result;
}

#endif
