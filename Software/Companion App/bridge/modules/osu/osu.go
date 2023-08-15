package osu

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"os/signal"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/utils"
	"runtime"

	"git.sr.ht/~sbinet/gg"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"golang.org/x/exp/slices"
)

var firstFrame = true
var keysStat = map[uint16]bool{}
var keysBrightness = map[uint16]float64{}

// https://xart3mis.github.io/posts/keylogger/
func keyTrackingThread(receptChan chan struct {
	uint16
	bool
}) {
	keyboardChan := make(chan types.KeyboardEvent, 1024)

	err := keyboard.Install(nil, keyboardChan)
	utils.CheckError(err)

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("start capturing keyboard input")

	for {
		select {
		case <-signalChan:
			fmt.Println("Received shutdown signal")
		case k := <-keyboardChan:
			if slices.Contains(config.Config.OsuTrackedKeys, uint16(k.VKCode)) {
				if k.Message == types.WM_KEYDOWN {
					receptChan <- struct {
						uint16
						bool
					}{uint16(k.VKCode), true}
				}
				if k.Message == types.WM_KEYUP {
					receptChan <- struct {
						uint16
						bool
					}{uint16(k.VKCode), false}
				}
			}
		}
	}
}

func startKeyTrackingThreads() {
	receptChan := make(chan struct {
		uint16
		bool
	})
	if runtime.GOOS == "windows" {
		go keyTrackingThread(receptChan)
	}
	for {
		pair := <-receptChan
		keysStat[pair.uint16] = pair.bool
	}
}

var OsuModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	if firstFrame {
		firstFrame = false
		go startKeyTrackingThreads()
	}

	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.MedSmallFont)
	anchorX := (float64(config.Config.CanvasRenderDimensions.X) - float64(len(config.Config.OsuTrackedKeys))*float64(config.Config.OsuTrackedKeysDimensions.X)) / 2
	anchorY := (float64(config.Config.CanvasRenderDimensions.Y) - float64(config.Config.OsuTrackedKeysDimensions.Y)) / 2
	for _, keycode := range config.Config.OsuTrackedKeys {
		if keysStat[keycode] {
			keysBrightness[keycode] = 1
		} else {
			keysBrightness[keycode] *= 0.6
		}
	}
	for i, keycode := range config.Config.OsuTrackedKeys {
		dc.SetRGB(keysBrightness[keycode], keysBrightness[keycode], keysBrightness[keycode])
		dc.DrawRectangle(anchorX+float64(i*config.Config.OsuTrackedKeysDimensions.X), anchorY, float64(config.Config.OsuTrackedKeysDimensions.X), float64(config.Config.OsuTrackedKeysDimensions.Y))
		dc.Fill()
		dc.SetRGB(1-keysBrightness[keycode], 1-keysBrightness[keycode], 1-keysBrightness[keycode])
		dc.DrawStringAnchored(config.Config.OsuTrackedKeysLabels[i], anchorX+float64(i*config.Config.OsuTrackedKeysDimensions.X)+float64(config.Config.OsuTrackedKeysDimensions.X)/2, float64(config.Config.CanvasRenderDimensions.Y)/2, 0.5, 0.5)
	}
	if !config.Config.RGBXMit {
		palette := []color.Color{
			color.Black,
			color.White,
		}
		d := dither.NewDitherer(palette)
		//d.Mapper = dither.Bayer(2, 2, 1.0)
		d.Matrix = dither.FloydSteinberg
		im = d.Dither(im).(*image.RGBA)
	}

	for i, _ := range config.Config.OsuTrackedKeys {
		dc.SetRGB(1, 1, 1)
		dc.DrawRectangle(anchorX+float64(i*config.Config.OsuTrackedKeysDimensions.X), anchorY, float64(config.Config.OsuTrackedKeysDimensions.X), float64(config.Config.OsuTrackedKeysDimensions.Y))
		dc.Stroke()
	}

	return renderer.AddWallpaperToFrame(renderer.RemoveAntiAliasing(im))
}}
