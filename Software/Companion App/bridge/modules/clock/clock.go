package clock

import (
	"image"
	"pscreenapp/bridge/modules"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"time"

	"github.com/fogleman/gg"
)

var ClockModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	now := time.Now()
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.BigFont)
	dc.DrawStringAnchored(now.Format("3:04:05 PM"), float64(config.CanvasRenderDimensions.X-4), float64(config.CanvasRenderDimensions.Y-4), 1, 0)
	dc.SetFontFace(renderer.MediumFont)
	dc.DrawStringAnchored(now.Format("1/2/2006"), 4, -4, 0, 1)
	return renderer.AddWallpaperToFrame(renderer.RemoveAntiAliasing(im))
}}
