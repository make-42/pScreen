pacman -S mingw-w64-x86_64-gcc
pacman -S mingw-w64-x86_64-pkg-config
pacman -S mingw-w64-x86_64-portaudio
pacman -S mingw-w64-x86_64-go
export GOROOT=/mingw64/lib/go
GOOS=windows CGO_ENABLED=1 CC=mingw-w64-x86_64-gcc go build
