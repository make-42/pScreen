package screensaver

import (
	"image"
	"pscreenapp/bridge/modules"
)

var ScreensaverModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	//fmt.Println("Rendering screensaver!")
	return im
}}
