package modules

import "image"

type Module struct {
	RenderFunction func(*image.RGBA) *image.RGBA
}
