package clock

import (
	"image"
	"math"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"time"

	"git.sr.ht/~sbinet/gg"
)

var ClockModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	now := time.Now()
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.BigFont)

	if config.Config.CircularScreenLayout {
		dc.DrawStringAnchored(now.Format("3:04 PM"), float64(config.Config.CanvasRenderDimensions.X)/2, float64(config.Config.CanvasRenderDimensions.Y)/2, 0.5, 0)
		dc.SetFontFace(renderer.MediumFont)
		dc.DrawStringAnchored(now.Format("1/2/2006"), float64(config.Config.CanvasRenderDimensions.X)/2, float64(config.Config.CanvasRenderDimensions.Y)/2, 0.5, 1)
		dc.DrawArc(float64(config.Config.CanvasRenderDimensions.X)/2, float64(config.Config.CanvasRenderDimensions.Y)/2, math.Max(float64(config.Config.CanvasRenderDimensions.Y), float64(config.Config.CanvasRenderDimensions.X))/2-4, 0, math.Pi*2*float64(now.Second())/60)
		dc.Stroke()
		return renderer.RemoveAntiAliasing(im)
	}
	dc.DrawStringAnchored(now.Format("3:04:05 PM"), float64(config.Config.CanvasRenderDimensions.X-4), float64(config.Config.CanvasRenderDimensions.Y-4), 1, 0)
	dc.SetFontFace(renderer.MediumFont)
	dc.DrawStringAnchored(now.Format("1/2/2006"), 4, -8, 0, 1)

	return renderer.AddWallpaperToFrame(renderer.RemoveAntiAliasing(im))
}}
