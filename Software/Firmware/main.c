/**
 * Copyright (c) 2020 Raspberry Pi (Trading) Ltd.
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "pico/stdlib.h"

int main() {
    stdio_init_all();
    while (true) {
        sleep_ms(250);
        printf("Blink!\n");
    }
}