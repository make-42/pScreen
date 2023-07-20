#include "zlib.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

#define FRAMEBUFFER_SIZE 2048
void get_frame(char * frame_buffer_output) {
  char frame_length_buffer[2];
  int i;
  unsigned short frame_length;
  // Receive screen update length
  for (i = 0; i < 2; ++i) {
    frame_length_buffer[i] = getchar();
  }
  memcpy(&frame_length, frame_length_buffer, 2);
  char compressed_frame_data[frame_length];
  // Receive compressed screen update
  for (i = 0; i < frame_length; ++i) {
    compressed_frame_data[i] = getchar();
  }
  // Decompress
  // zlib struct
  z_stream infstream;
  infstream.zalloc = Z_NULL;
  infstream.zfree = Z_NULL;
  infstream.opaque = Z_NULL;
  // setup "b" as the input and "c" as the compressed output
  infstream.avail_in = (uInt)frame_length;            // size of input
  infstream.next_in = (Bytef *)compressed_frame_data; // input char array
  infstream.avail_out = (uInt)FRAMEBUFFER_SIZE;       // size of output
  infstream.next_out = (Bytef *)frame_buffer_output;         // output char array

  // the actual DE-compression work.
  inflateInit(&infstream);
  inflate(&infstream, Z_NO_FLUSH);
  inflateEnd(&infstream);
  /*for (i = 0; i < FRAMEBUFFER_SIZE; i++) {
    printf("%x", frame_buffer[i]);
  }
  printf("\n");*/
};

void signal_ready_to_receive_next_frame() {
  //  Send unblock command
  printf("O");
}