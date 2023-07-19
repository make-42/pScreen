cmake_minimum_required(VERSION 3.13)

# initialize the SDK based on PICO_SDK_PATH
# note: this must happen before project()
include(pico_sdk_import.cmake)

project(pscreen C CXX ASM)

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -ffunction-sections -fdata-sections -Wl,--gc-sections")
set(CMAKE_C_STANDARD 11)
set(CMAKE_CXX_STANDARD 17)
set(PICO_BOARD pico)

# initialize the Raspberry Pi Pico SDK
pico_sdk_init()

# rest of your project

add_executable(pscreen-ssd1306
	ssd1306ver.c pico-ssd1306/ssd1306.c
)

add_executable(pscreen-ssd1322
	ssd1322ver.c
)

target_include_directories(pscreen-ssd1306
    PUBLIC
        ${CMAKE_CURRENT_LIST_DIR}/../
)

target_include_directories(pscreen-ssd1322
    PUBLIC
        ${CMAKE_CURRENT_LIST_DIR}/../
)
target_include_directories(pscreen-ssd1322
    PRIVATE
        u8g2/csrc
)

file(GLOB U8G2_SRC u8g2/csrc/*.c)
add_library(u8g2 ${U8G2_SRC})

target_link_options(pscreen-ssd1306 PRIVATE -Xlinker --print-memory-usage)
target_compile_options(pscreen-ssd1306 PRIVATE -Wall -Wextra -DCFG_TUSB_DEBUG=1)

target_link_options(pscreen-ssd1322 PRIVATE -Xlinker --print-memory-usage)
target_compile_options(pscreen-ssd1322 PRIVATE -Wall -Wextra -DCFG_TUSB_DEBUG=1)


target_link_libraries(pscreen-ssd1306 pico_stdlib hardware_i2c)
target_link_libraries(pscreen-ssd1322 pico_stdlib u8g2 hardware_spi)

pico_enable_stdio_usb(pscreen-ssd1306 1)
pico_enable_stdio_uart(pscreen-ssd1306 0)
pico_enable_stdio_usb(pscreen-ssd1322 1)
pico_enable_stdio_uart(pscreen-ssd1322 0)

# create map/bin/hex/uf2 file in addition to ELF.
pico_add_extra_outputs(pscreen-ssd1306)
pico_add_extra_outputs(pscreen-ssd1322)