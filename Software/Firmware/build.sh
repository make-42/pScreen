mkdir build
cd build
export PICO_SDK_PATH=../pico-sdk/
echo $PICO_SDK_PATH
cmake ..
make -j$(nproc)
cd ..
