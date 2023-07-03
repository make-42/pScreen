package encoder

import (
	"image"
	"pscreenapp/config"
	"pscreenapp/utils"
)

func EncodeFrameToBytes(im *image.RGBA) []byte {
	var bytes []byte
	for x := 0; x < config.CanvasRenderDimensions[0]; x++ {
		for y := 0; y < config.CanvasRenderDimensions[1]/8; y++ {
			currentByte := byte(0)
			for z := 0; z < 8; z++ {
				rgba := im.RGBAAt(x, y*8+z)
				if rgba.R == 255 {
					currentByte += byte(utils.IntPow(2, 7-z))
				}
			}
			bytes = append(bytes, currentByte)
		}
	}
	return bytes
}
