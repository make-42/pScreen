package clock

import (
	"image"
	"image/draw"
	"pscreenapp/bridge/modules"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"time"

	"github.com/fogleman/gg"
)

var ClockModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	now := time.Now()
	//fmt.Println("Rendering clock!")
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.BigFont)
	dc.DrawStringAnchored(now.Format("15:04:05"), float64(config.CanvasRenderDimensions[0]-4), float64(config.CanvasRenderDimensions[1]-4), 1, 0)
	dc.SetFontFace(renderer.MediumFont)
	dc.DrawStringAnchored(now.Format("1/2/2006"), 4, -4, 0, 1)
	b := im.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), renderer.BackgroundImage.(*image.RGBA), b.Min, draw.Src)

	return renderer.CompositeBackgroundAndForeground(m, renderer.RemoveAntiAliasing(im))
}}
