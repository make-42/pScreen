package blank

import (
	"image"
	"pscreenapp/bridge/modules"
)

var BlankModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	return im
}}
