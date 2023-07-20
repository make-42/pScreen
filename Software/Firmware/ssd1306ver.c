#include "pico-ssd1306/ssd1306.h"
#include "pico/stdlib.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "comms/comms.h"

#define CANVAS_WIDTH 256
#define CANVAS_HEIGHT 64

void setup_gpios(void);

void setup_gpios(void) {
  i2c_init(i2c1, 400000);
  gpio_set_function(6, GPIO_FUNC_I2C);
  gpio_set_function(7, GPIO_FUNC_I2C);
  gpio_pull_up(6);
  gpio_pull_up(7);

  i2c_init(i2c0, 400000);
  gpio_set_function(12, GPIO_FUNC_I2C);
  gpio_set_function(13, GPIO_FUNC_I2C);
  gpio_pull_up(12);
  gpio_pull_up(13);
}

int main() {
  stdio_init_all();
  printf("configuring pins...\n");
  setup_gpios();
  int i;
  int x;
  int y;
  int z;
  char frame_buffer[FRAMEBUFFER_SIZE];
  char *frame_buffer_address;
  ssd1306_t disp_a;
  ssd1306_t disp_b;
  disp_a.external_vcc = false;
  disp_b.external_vcc = false;
  ssd1306_init(&disp_a, 128, 64, 0x3C, i2c1);
  ssd1306_clear(&disp_a);
  ssd1306_init(&disp_b, 128, 64, 0x3C, i2c0);
  ssd1306_clear(&disp_b);
  ssd1306_draw_string(&disp_a, 0, 0, 1, "Waiting for");
  ssd1306_draw_string(&disp_b, 0, 0, 1, "companion app...");
  ssd1306_show(&disp_a);
  ssd1306_show(&disp_b);
  //  Main loop
  for (;;) {
    // Receive screen update
    frame_buffer_address = get_frame();
    memcpy(&frame_buffer, frame_buffer_address, FRAMEBUFFER_SIZE);
    free(frame_buffer_address);
    // Screen draw logic
    ssd1306_clear(&disp_a);
    ssd1306_clear(&disp_b);
    for (x = 0; x < CANVAS_WIDTH; x++) {
      for (y = 0; y < CANVAS_HEIGHT / 8; y++) {
        for (z = 0; z < 8; z++) {
          if ((frame_buffer[x * CANVAS_HEIGHT / 8 + y] >> (7 - z)) & 0x01) {
            if (x > 127) {
              ssd1306_draw_pixel(&disp_b, x - 128, y * 8 + z);
            } else {
              ssd1306_draw_pixel(&disp_a, x, y * 8 + z);
            }
          }
        }
      }
    }
    ssd1306_show(&disp_a);
    ssd1306_show(&disp_b);
    //  Send unblock command
    signal_ready_to_receive_next_frame();
  }
}