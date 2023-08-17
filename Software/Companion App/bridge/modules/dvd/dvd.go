package dvd

import (
	"image"
	"image/draw"
	"math"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/utils"
)

var firstFrame = true
var DVDPosition utils.CoordsFloat
var DVDVelocity utils.CoordsFloat

var DVDModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	if firstFrame {
		firstFrame = false
		DVDPosition = utils.CoordsFloat{X: float64(config.Config.CanvasRenderDimensions.X-renderer.DVDLogo.Bounds().Dx()) / 2, Y: float64(config.Config.CanvasRenderDimensions.Y-renderer.DVDLogo.Bounds().Dy()) / 2}
		angle := 2 * math.Pi / 360 * 45

		DVDVelocity = utils.CoordsFloat{
			X: math.Cos(angle) * config.Config.DVDLogoVelocity,
			Y: -math.Sin(angle) * config.Config.DVDLogoVelocity,
		}
	}
	DVDPosition = utils.CoordsFloat{DVDPosition.X + DVDVelocity.X, DVDPosition.Y + DVDVelocity.Y}

	logoB := renderer.DVDLogo.Bounds()
	if DVDPosition.X+float64(logoB.Dx()) > float64(config.Config.CanvasRenderDimensions.X) {
		DVDPosition.X = float64(config.Config.CanvasRenderDimensions.X - logoB.Dx())
		DVDVelocity.X = -DVDVelocity.X
	}
	if DVDPosition.X < 0 {
		DVDVelocity.X = -DVDVelocity.X
		DVDPosition.X = 0
	}
	if DVDPosition.Y+float64(logoB.Dy()) > float64(config.Config.CanvasRenderDimensions.Y) {
		DVDPosition.Y = float64(config.Config.CanvasRenderDimensions.Y - logoB.Dy())
		DVDVelocity.Y = -DVDVelocity.Y
	}
	if DVDPosition.Y < 0 {
		DVDVelocity.Y = -DVDVelocity.Y
		DVDPosition.Y = 0
	}

	draw.Draw(im, logoB.Add(image.Pt(int(DVDPosition.X), int(DVDPosition.Y))), renderer.DVDLogo, logoB.Min, draw.Src)

	return renderer.RemoveAntiAliasing(im)
}}
