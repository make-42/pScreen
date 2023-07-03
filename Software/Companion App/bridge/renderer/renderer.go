package renderer

import (
	"image"
	"image/draw"
	"pscreenapp/bridge/encoder"
	"pscreenapp/bridge/modules"
	"pscreenapp/config"
	"pscreenapp/utils"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

var BigFont font.Face
var MediumFont font.Face

func LoadRendererSharedRessources() {
	var err error
	BigFont, err = gg.LoadFontFace("./assets/fonts/BegikaFixed.ttf", 40)
	utils.CheckError(err)
	MediumFont, err = gg.LoadFontFace("./assets/fonts/BegikaFixed.ttf", 30)
	utils.CheckError(err)
}

func RenderFrame(mod modules.Module) []byte {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{config.CanvasRenderDimensions[0], config.CanvasRenderDimensions[1]}
	im := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	im = mod.RenderFunction(im)
	//f, _ := os.Create("image.png")
	//png.Encode(f, im)
	return encoder.EncodeFrameToBytes(im)

}

func RemoveAntiAliasing(im *image.RGBA) *image.RGBA {
	img := imaging.AdjustContrast(im, 127)
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m
}
