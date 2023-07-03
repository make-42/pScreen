package media

import (
	"image"
	"pscreenapp/bridge/modules"
)

var MediaModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	//fmt.Println("Rendering media!")
	return im
}}
