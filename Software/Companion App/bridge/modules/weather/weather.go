package weather

import (
	"image"
	"pscreenapp/bridge/modules"
)

var WeatherModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	//fmt.Println("Rendering weather!")
	return im
}}
