package elements

import (
	"math"
	"pscreenapp/config"
	"pscreenapp/utils"

	"github.com/fogleman/gg"
)

func GetMediaProgressBarYFromX(x float64) float64 {
	return float64(config.CanvasRenderDimensions.Y-4) - float64(config.MediaProgressBarIndicatorRadius) + (math.Sin(float64(x)*math.Pi*config.MediaProgressBarWaveScale)-1)*(float64(config.MediaProgressBarHeight)/2-2)
}

func DrawMediaProgressBar(dc *gg.Context, positionSec, durationSec float64) {
	// This function does not handle the font
	dc.DrawStringAnchored(utils.FormatDuration(positionSec), 4, float64(config.CanvasRenderDimensions.Y-8)-float64(config.MediaProgressBarHeight), 0, 0)
	dc.DrawStringAnchored(utils.FormatDuration(durationSec), float64(config.CanvasRenderDimensions.X-4), float64(config.CanvasRenderDimensions.Y-8)-float64(config.MediaProgressBarHeight), 1, 0)
	mediaProgressBarWidth := config.CanvasRenderDimensions.X - 8
	for x := 0; x < mediaProgressBarWidth; x++ {
		dc.DrawCircle(float64(x+4), GetMediaProgressBarYFromX(float64(x)), 1)
		dc.Fill()
	}
	currentMediaPositionX := 0.0
	if durationSec != 0 {
		currentMediaPositionX = float64(mediaProgressBarWidth) * positionSec / durationSec
	}
	dc.DrawCircle(4+currentMediaPositionX, GetMediaProgressBarYFromX(currentMediaPositionX), float64(config.MediaProgressBarIndicatorRadius))
	dc.Fill()
}
