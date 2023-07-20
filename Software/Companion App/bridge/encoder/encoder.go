package encoder

import (
	"encoding/binary"
	"image"
	"pscreenapp/config"
	"pscreenapp/utils"

	"bytes"
	"compress/zlib"
)

var UncompressedBytesN = -1
var CompressedBytesN = -1

func EncodeFrameToBytes(im *image.RGBA) []byte {
	var uncompressedBytes []byte
	for x := 0; x < config.CanvasRenderDimensions.X; x++ {
		for y := 0; y < config.CanvasRenderDimensions.Y/8; y++ {
			currentByte := byte(0)
			for z := 0; z < 8; z++ {
				rgba := im.RGBAAt(x, y*8+z)
				if rgba.R == 255 {
					currentByte += byte(utils.IntPow(2, 7-z))
				}
			}
			uncompressedBytes = append(uncompressedBytes, currentByte)
		}
	}
	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, zlib.BestCompression) // requires no writer if WriteBuffer is used
	utils.CheckError(err)
	_, err = w.Write([]byte(uncompressedBytes))
	utils.CheckError(err)
	w.Close()
	compressedBytes := b.Bytes()
	outputBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(outputBytes, uint16(len(compressedBytes)))
	outputBytes = append(outputBytes, compressedBytes...)
	UncompressedBytesN = len(uncompressedBytes)
	CompressedBytesN = len(outputBytes)
	return outputBytes
}
