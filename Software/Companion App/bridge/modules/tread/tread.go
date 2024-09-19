package tread

import (
	"fmt"
	"image"
	"math"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/fogleman/ease"
	"golang.org/x/image/font"
)

func DrawBeltHorizontal(start, end int, spacing, width int, fontface font.Face) image.Image {
	length := spacing * (end - start + 1) * 2
	dc := gg.NewContext(length, width)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(fontface)
	for i := 0; i < (end-start+1)*2; i++ {
		dc.DrawStringAnchored(fmt.Sprintf("%d", start+(i)%(end-start+1)), float64(spacing)*(float64(i)+0.5), float64(width)/2, 0.5, 0.25)
	}
	dc.DrawLine(0, 0, float64(length), 0)
	dc.Stroke()
	dc.DrawLine(0, float64(width)-0.5, float64(length), float64(width)-0.5)
	dc.Stroke()
	return dc.Image()
}

func DrawBeltDiagonal(start, end int, spacing, width int, fontface font.Face) image.Image {
	length := spacing * (end - start + 1) * 2
	dc := gg.NewContext(width+length, length)
	dc.SetRGB(0, 0, 0)
	dc.MoveTo(float64(length+width), 0)
	dc.LineTo(float64(width), float64(length)-0.5)
	dc.LineTo(0, float64(length)-0.5)
	dc.LineTo(float64(length), 0)
	dc.ClosePath()
	dc.FillPreserve()
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.Stroke()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(fontface)
	for i := 0; i < (end-start+1)*2; i++ {
		dc.DrawStringAnchored(fmt.Sprintf("%2d", start+(i)%(end-start+1)), float64(width)/2+float64(spacing)*(float64(i)+0.5), float64(length)-float64(spacing)*(float64(i)+0.5), 0.5, 0.1)
	}
	return dc.Image()
}

var firstFrame = true
var HourBelt image.Image
var MinuteABelt image.Image
var MinuteBBelt image.Image

var TreadModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	if firstFrame {
		HourBelt = DrawBeltHorizontal(1, 12, 32, 32, renderer.MedSmallFont)
		MinuteABelt = DrawBeltDiagonal(0, 5, 32, 24, renderer.MedSmallerFont)
		MinuteBBelt = DrawBeltDiagonal(0, 9, 32, 24, renderer.MedSmallerFont)
		firstFrame = false
	}
	now := time.Now()
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	aoffset := -((now.Second()/10+4)%6-3)*32 + 16 - int(ease.InOutQuint((math.Max(float64(now.UnixMilli()%10000), 9600)-9600)/400)*32)
	boffset := -((now.Second()%10+6)%10-5)*32 + 16 - int(ease.InOutQuint((math.Max(float64(now.UnixMilli()%1000), 600)-600)/400)*32)
	coffset := -((now.Minute()/10+4)%6-3)*32 + 16 - int(ease.InOutQuint((math.Max(float64(now.UnixMilli()%600000), 599600)-599600)/400)*32)
	doffset := -((now.Minute()%10+6)%10-5)*32 + 16 - int(ease.InOutQuint((math.Max(float64(now.UnixMilli()%60000), 59600)-59600)/400)*32)
	eoffset := -((now.Hour()+6)%12-6)*32 - int(ease.InOutQuint((math.Max(float64(now.UnixMilli()%3600000), 3599600)-3599600)/400)*32)

	dc.DrawImageAnchored(HourBelt, int(float64(config.Config.CanvasRenderDimensions.X)/2)+16+eoffset-80, int(float64(config.Config.CanvasRenderDimensions.Y)/2), 0.5, 0.5)

	dc.DrawImageAnchored(MinuteABelt, int(float64(config.Config.CanvasRenderDimensions.X)/2)+coffset-58+16, int(float64(config.Config.CanvasRenderDimensions.Y)/2)-coffset+4, 0.5, 0.5)
	dc.DrawImageAnchored(MinuteBBelt, int(float64(config.Config.CanvasRenderDimensions.X)/2)+doffset-32+16, int(float64(config.Config.CanvasRenderDimensions.Y)/2)-doffset+4, 0.5, 0.5)
	dc.DrawImageAnchored(MinuteABelt, int(float64(config.Config.CanvasRenderDimensions.X)/2)+aoffset+32-32+16, int(float64(config.Config.CanvasRenderDimensions.Y)/2)-aoffset+4, 0.5, 0.5)
	dc.DrawImageAnchored(MinuteBBelt, int(float64(config.Config.CanvasRenderDimensions.X)/2)+boffset+58-32+16, int(float64(config.Config.CanvasRenderDimensions.Y)/2)-boffset+4, 0.5, 0.5)

	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(2.0)
	dc.DrawRectangle(float64(config.Config.CanvasRenderDimensions.X)/2-16-80, float64(config.Config.CanvasRenderDimensions.Y)/2-16, 32, 32)
	dc.Stroke()
	dc.DrawRectangle(float64(config.Config.CanvasRenderDimensions.X)/2-48, float64(config.Config.CanvasRenderDimensions.Y)/2-12, 48, 24)
	dc.Stroke()
	dc.DrawRectangle(float64(config.Config.CanvasRenderDimensions.X)/2+10, float64(config.Config.CanvasRenderDimensions.Y)/2-12, 48, 24)
	dc.Stroke()

	dc.SetFontFace(renderer.TinyFont)
	dc.SetRGB(1, 1, 1)
	dc.DrawRoundedRectangle(float64(config.Config.CanvasRenderDimensions.X)-16-4, 4, 16, 10, 2)
	dc.Fill()
	dc.DrawRoundedRectangle(float64(config.Config.CanvasRenderDimensions.X)-64-4, float64(config.Config.CanvasRenderDimensions.Y)-10-4, 64, 10, 2)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringWrapped(now.Format("PM"), float64(config.Config.CanvasRenderDimensions.X)-16+4, 6, 0.5, 0.5, 16, 0, gg.AlignCenter)
	dc.DrawStringWrapped(now.Format("01/02/2006"), float64(config.Config.CanvasRenderDimensions.X)-64/2-4, float64(config.Config.CanvasRenderDimensions.Y)-12, 0.5, 0.5, 64, 0, gg.AlignCenter)

	return renderer.RemoveAntiAliasing(im)
}}
