package monitor

import (
	"fmt"
	"image"
	"pscreenapp/bridge/modules"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"pscreenapp/utils"
	"time"

	"github.com/fogleman/gg"
	"github.com/shirou/gopsutil/cpu"
)

var MonitorModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	cpuInfo, err := cpu.Info()
	utils.CheckError(err)
	cpuPercent, err := cpu.Percent(time.Millisecond*config.CPUUsageSamplingMilliseconds, true)
	utils.CheckError(err)
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.TinyFont)
	dc.DrawStringAnchored(cpuInfo[0].ModelName, 4, 2, 0, 1)
	percentSum := 0.0
	for i, corePercent := range cpuPercent {
		rectangleX := 4 + i*(config.CPUUsageBarDimensions[0]+config.CPUUsageBarMargin)
		rectangleY := config.CanvasRenderDimensions[1] - config.CPUUsageBarDimensions[1] - 4
		rectanglePercentageHeight := float64(config.CPUUsageBarDimensions[1]-1)*corePercent/100 + 1
		dc.DrawRectangle(float64(rectangleX), float64(rectangleY+config.CPUUsageBarDimensions[1])-rectanglePercentageHeight, float64(config.CPUUsageBarDimensions[0]), float64(rectanglePercentageHeight))
		dc.DrawRectangle(float64(rectangleX), float64(rectangleY), float64(config.CPUUsageBarDimensions[0]), 1)
		dc.Fill()
		percentSum += corePercent
	}
	dc.SetFontFace(renderer.MediumFont)
	dc.DrawStringAnchored(fmt.Sprintf("%0.1f %%", percentSum/float64(len(cpuPercent))), float64(config.CanvasRenderDimensions[0]-4), float64(config.CanvasRenderDimensions[1]-4), 1, 0)

	return renderer.RemoveAntiAliasing(im)
}}
