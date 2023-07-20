#include "hardware/gpio.h"
#include "hardware/spi.h"
#include "pico/stdlib.h"
#include <stdio.h>
#include <string.h>
#include <u8g2.h>

#define BLOCKSIZE 2048
#define CANVAS_WIDTH 256
#define CANVAS_HEIGHT 64

u8g2_t u8g2;

// Display settings
#define SPI_PORT spi0
#define PIN_CS 21
#define PIN_SCK 18
#define PIN_MOSI 19
#define SPI_SPEED 4000 * 1000
#define PIN_DC 15
#define PIN_RST 14

uint8_t u8x8_byte_pico_hw_spi(u8x8_t *u8x8, uint8_t msg, uint8_t arg_int,
                              void *arg_ptr) {
  uint8_t *data;
  switch (msg) {
  case U8X8_MSG_BYTE_SEND:
    data = (uint8_t *)arg_ptr;
    spi_write_blocking(SPI_PORT, data, arg_int);
    break;
  case U8X8_MSG_BYTE_INIT:
    u8x8_gpio_SetCS(u8x8, u8x8->display_info->chip_disable_level);
    break;
  case U8X8_MSG_BYTE_SET_DC:
    u8x8_gpio_SetDC(u8x8, arg_int);
    break;
  case U8X8_MSG_BYTE_START_TRANSFER:
    u8x8_gpio_SetCS(u8x8, u8x8->display_info->chip_enable_level);
    u8x8->gpio_and_delay_cb(u8x8, U8X8_MSG_DELAY_NANO,
                            u8x8->display_info->post_chip_enable_wait_ns, NULL);
    break;
  case U8X8_MSG_BYTE_END_TRANSFER:
    u8x8->gpio_and_delay_cb(u8x8, U8X8_MSG_DELAY_NANO,
                            u8x8->display_info->pre_chip_disable_wait_ns, NULL);
    u8x8_gpio_SetCS(u8x8, u8x8->display_info->chip_disable_level);
    break;
  default:
    return 0;
  }
  return 1;
}

uint8_t u8x8_gpio_and_delay_pico(u8x8_t *u8x8, uint8_t msg, uint8_t arg_int,
                                 void *arg_ptr) {
  switch (msg) {
  case U8X8_MSG_GPIO_AND_DELAY_INIT:
    spi_init(SPI_PORT, SPI_SPEED);
    gpio_set_function(PIN_CS, GPIO_FUNC_SIO);
    gpio_set_function(PIN_SCK, GPIO_FUNC_SPI);
    gpio_set_function(PIN_MOSI, GPIO_FUNC_SPI);
    gpio_init(PIN_RST);
    gpio_init(PIN_DC);
    gpio_init(PIN_CS);
    gpio_set_dir(PIN_RST, GPIO_OUT);
    gpio_set_dir(PIN_DC, GPIO_OUT);
    gpio_set_dir(PIN_CS, GPIO_OUT);
    gpio_put(PIN_RST, 1);
    gpio_put(PIN_CS, 1);
    gpio_put(PIN_DC, 0);
    break;
  case U8X8_MSG_DELAY_NANO: // delay arg_int * 1 nano second
    sleep_us(arg_int); // 1000 times slower, though generally fine in practice
                       // given rp2040 has no `sleep_ns()`
    break;
  case U8X8_MSG_DELAY_100NANO: // delay arg_int * 100 nano seconds
    sleep_us(arg_int);
    break;
  case U8X8_MSG_DELAY_10MICRO: // delay arg_int * 10 micro seconds
    sleep_us(arg_int * 10);
    break;
  case U8X8_MSG_DELAY_MILLI: // delay arg_int * 1 milli second
    sleep_ms(arg_int);
    break;
  case U8X8_MSG_GPIO_CS: // CS (chip select) pin: Output level in arg_int
    gpio_put(PIN_CS, arg_int);
    break;
  case U8X8_MSG_GPIO_DC: // DC (data/cmd, A0, register select) pin: Output level
    gpio_put(PIN_DC, arg_int);
    break;
  case U8X8_MSG_GPIO_RESET:     // Reset pin: Output level in arg_int
    gpio_put(PIN_RST, arg_int); // printf("U8X8_MSG_GPIO_RESET %d\n", arg_int);
    break;
  default:
    u8x8_SetGPIOResult(u8x8, 1); // default return value
    break;
  }
  return 1;
}

void display_init() {
  u8g2_Setup_ssd1322_nhd_256x64_f(
      &u8g2, U8G2_R2, u8x8_byte_pico_hw_spi,
      u8x8_gpio_and_delay_pico); // init u8g2 structure
  u8g2_InitDisplay(&u8g2);     // send init sequence to the display, display is
                               // in sleep mode after this,
  u8g2_SetPowerSave(&u8g2, 0); // wake up display
}

void waiting_for_companion_app_message() {
  u8g2_ClearBuffer(&u8g2);
  u8g2_SetDrawColor(&u8g2, 1);
  u8g2_SetFont(&u8g2, u8g2_font_t0_11_te);
  u8g2_DrawStr(&u8g2, 0, 20, "Waiting for companion app...");
  u8g2_UpdateDisplay(&u8g2);
}

int main() {
  stdio_init_all();
  display_init();
  waiting_for_companion_app_message();
  int i;
  int x;
  int y;
  int z;
  char frame_buffer[BLOCKSIZE];
  //  Main loop
  for (;;) {
    // Receive screen update
    for (i = 0; i < BLOCKSIZE; ++i) {
      frame_buffer[i] = getchar();
    }
    // Screen draw logic
    u8g2_ClearBuffer(&u8g2);
    u8g2_SetDrawColor(&u8g2, 1);

    for (x = 0; x < CANVAS_WIDTH; x++) {
      for (y = 0; y < CANVAS_HEIGHT / 8; y++) {
        for (z = 0; z < 8; z++) {
          if ((frame_buffer[x * CANVAS_HEIGHT / 8 + y] >> (7 - z)) & 0x01) {
            u8g2_DrawPixel(&u8g2, x, y * 8 + z);
          }
        }
      }
    }
    u8g2_UpdateDisplay(&u8g2);

    //  Send unblock command
    printf("O");
  }
}