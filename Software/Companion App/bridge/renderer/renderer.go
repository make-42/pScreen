package renderer

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
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
var MedSmallFont font.Face
var SmallFont font.Face
var TinyFont font.Face
var BackgroundImage image.Image

func LoadRendererSharedRessources() {
	var err error
	BigFont, err = gg.LoadFontFace("./assets/fonts/BegikaFixed.ttf", 40)
	utils.CheckError(err)
	MediumFont, err = gg.LoadFontFace("./assets/fonts/BegikaFixed.ttf", 30)
	utils.CheckError(err)
	MedSmallFont, err = gg.LoadFontFace("./assets/fonts/BegikaFixed.ttf", 24)
	utils.CheckError(err)
	SmallFont, err = gg.LoadFontFace("./assets/fonts/BegikaFixed.ttf", 16)
	utils.CheckError(err)
	TinyFont, err = gg.LoadFontFace("./assets/fonts/Lato-Bold.ttf", 12)
	utils.CheckError(err)
	imgFile, err := os.Open("./assets/img/bg.png")
	defer imgFile.Close()
	utils.CheckError(err)
	BackgroundImage, _, err = image.Decode(imgFile)
	utils.CheckError(err)
}

func RenderFrame(mod modules.Module) []byte {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{config.CanvasRenderDimensions.X, config.CanvasRenderDimensions.Y}
	im := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	im = mod.RenderFunction(im)
	if config.DebugSaveScreen {
		f, _ := os.Create("image.png")
		png.Encode(f, im)
	}
	return encoder.EncodeFrameToBytes(im)

}

func NRGBAImgToRGBAImg(im *image.NRGBA) *image.RGBA {
	b := im.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), im, b.Min, draw.Src)
	return m
}

func RemoveAntiAliasing(im *image.RGBA) *image.RGBA {
	return NRGBAImgToRGBAImg(imaging.AdjustContrast(im, 127))
}

func CompositeBackgroundAndForeground(bgImg *image.RGBA, fgImg *image.RGBA) *image.RGBA {
	fgBorder := NRGBAImgToRGBAImg(imaging.Invert(imaging.AdjustContrast(imaging.AdjustBrightness(imaging.Blur(fgImg, 1), 40), 127)))
	b := fgBorder.Bounds()
	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {
			if fgImg.RGBAAt(x, y).R == 255 {
				bgImg.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else if fgBorder.RGBAAt(x, y).R == 0 {
				bgImg.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return bgImg
}

func AddWallpaperToFrame(fgImg *image.RGBA) *image.RGBA {
	b := fgImg.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), BackgroundImage.(*image.RGBA), b.Min, draw.Src)
	return CompositeBackgroundAndForeground(m, fgImg)
}
