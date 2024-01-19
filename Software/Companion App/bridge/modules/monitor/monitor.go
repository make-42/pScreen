package monitor

import (
	"fmt"
	"image"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/utils"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/shirou/gopsutil/cpu"

	"github.com/ztrue/tracerr"
)

var MonitorModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	cpuInfo, err := cpu.Info()
	utils.CheckError(tracerr.Wrap(err))
	cpuPercent, err := cpu.Percent(time.Millisecond*time.Duration(config.Config.CPUUsageSamplingMilliseconds), true)
	utils.CheckError(tracerr.Wrap(err))
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.TinyFont)
	dc.DrawStringAnchored(cpuInfo[0].ModelName, 4, 2, 0, 1)
	percentSum := 0.0
	for i, corePercent := range cpuPercent {
		rectangleX := 4 + i*(config.Config.CPUUsageBarDimensions.X+config.Config.CPUUsageBarMargin)
		rectangleY := config.Config.CanvasRenderDimensions.Y - config.Config.CPUUsageBarDimensions.Y - 4
		rectanglePercentageHeight := float64(config.Config.CPUUsageBarDimensions.Y-1)*corePercent/100 + 1
		dc.DrawRectangle(float64(rectangleX), float64(rectangleY+config.Config.CPUUsageBarDimensions.Y)-rectanglePercentageHeight, float64(config.Config.CPUUsageBarDimensions.X), float64(rectanglePercentageHeight))
		dc.DrawRectangle(float64(rectangleX), float64(rectangleY), float64(config.Config.CPUUsageBarDimensions.X), 1)
		dc.Fill()
		percentSum += corePercent
	}
	dc.SetFontFace(renderer.MediumFont)
	dc.DrawStringAnchored(fmt.Sprintf("%0.1f %%", percentSum/float64(len(cpuPercent))), float64(config.Config.CanvasRenderDimensions.X-4), float64(config.Config.CanvasRenderDimensions.Y-4), 1, 0)

	return renderer.RemoveAntiAliasing(im)
}}
