package renderer

import (
	"embed"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"pscreenapp/bridge/encoder"
	"pscreenapp/bridge/modules"
	"pscreenapp/config"
	"pscreenapp/utils"
	"strings"

	"git.sr.ht/~sbinet/gg"
	"github.com/disintegration/imaging"
	"github.com/pbnjay/pixfont"
	"golang.org/x/image/font"
)

//go:embed assets/fonts/iosevka-heavy.ttf
//go:embed assets/fonts/iosevka-medium.ttf
//go:embed assets/img/bg.png
//go:embed assets/gif/komi-san-48.gif
var assetsFolder embed.FS

var BigFont font.Face
var MediumFont font.Face
var MedSmallFont font.Face
var MedSmallerFont font.Face
var SmallFont font.Face
var TinyFont font.Face
var BackgroundImage image.Image
var QRCodeModuleGIF *gif.GIF

func LoadRendererSharedRessources() {
	var err error
	BigFont, err = gg.LoadFontFaceFromFS(assetsFolder, "assets/fonts/iosevka-heavy.ttf", 40)
	utils.CheckError(err)
	MediumFont, err = gg.LoadFontFaceFromFS(assetsFolder, "assets/fonts/iosevka-heavy.ttf", 30)
	utils.CheckError(err)
	MedSmallFont, err = gg.LoadFontFaceFromFS(assetsFolder, "assets/fonts/iosevka-heavy.ttf", 24)
	utils.CheckError(err)
	MedSmallerFont, err = gg.LoadFontFaceFromFS(assetsFolder, "assets/fonts/iosevka-heavy.ttf", 20)
	utils.CheckError(err)
	SmallFont, err = gg.LoadFontFaceFromFS(assetsFolder, "assets/fonts/iosevka-heavy.ttf", 16) // BegikaFixed.ttf
	utils.CheckError(err)
	TinyFont, err = gg.LoadFontFaceFromFS(assetsFolder, "assets/fonts/iosevka-medium.ttf", 12) // Lato-Bold.ttf
	utils.CheckError(err)
	imgFile, err := assetsFolder.Open("assets/img/bg.png")
	defer imgFile.Close()
	utils.CheckError(err)
	BackgroundImage, _, err = image.Decode(imgFile)
	utils.CheckError(err)
	gifFile, err := assetsFolder.Open("assets/gif/komi-san-48.gif")
	defer gifFile.Close()
	utils.CheckError(err)
	QRCodeModuleGIF, err = gif.DecodeAll(gifFile)
	utils.CheckError(err)
}

func RenderFrame(mod modules.Module) []byte {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{config.Config.CanvasRenderDimensions.X, config.Config.CanvasRenderDimensions.Y}
	im := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	im = mod.RenderFunction(im)
	if config.Config.DebugSaveScreen {
		f, _ := os.Create("image.png")
		png.Encode(f, im)
	}
	if config.Config.RotateScreen180Deg {
		im = RotateImage180deg(im)
	}
	return encoder.EncodeFrameToBytes(im)

}

func RotateImage180deg(im *image.RGBA) *image.RGBA {
	return NRGBAImgToRGBAImg(imaging.Rotate180(im))
}

func NRGBAImgToRGBAImg(im *image.NRGBA) *image.RGBA {
	b := im.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), im, b.Min, draw.Src)
	return m
}

func YCbCrImgToRGBAImg(im *image.YCbCr) *image.RGBA {
	b := im.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), im, b.Min, draw.Src)
	return m
}

func RemoveAntiAliasing(im *image.RGBA) *image.RGBA {
	if config.Config.RGBXMit {
		return im
	}
	return NRGBAImgToRGBAImg(imaging.AdjustContrast(im, 127))
}

func CompositeBackgroundAndForeground(bgImg *image.RGBA, fgImg *image.RGBA) *image.RGBA {
	ba := bgImg.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, ba.Dx(), ba.Dy()))
	draw.Draw(m, m.Bounds(), bgImg, ba.Min, draw.Src)
	fgBorder := NRGBAImgToRGBAImg(imaging.Invert(imaging.AdjustContrast(imaging.AdjustBrightness(imaging.Blur(fgImg, 1), 40), 127)))
	b := fgBorder.Bounds()
	if config.Config.RGBXMit {
		for x := 0; x < b.Dx(); x++ {
			for y := 0; y < b.Dy(); y++ {
				if (fgBorder.RGBAAt(x, y) != color.RGBA{255, 255, 255, 255}) {
					m.Set(x, y, fgBorder.RGBAAt(x, y))
				}
			}
		}
		for x := 0; x < b.Dx(); x++ {
			for y := 0; y < b.Dy(); y++ {
				if (fgImg.RGBAAt(x, y) != color.RGBA{0, 0, 0, 255}) {
					m.Set(x, y, fgImg.RGBAAt(x, y))
				}
			}
		}
	} else {
		for x := 0; x < b.Dx(); x++ {
			for y := 0; y < b.Dy(); y++ {
				if fgImg.RGBAAt(x, y).R == 255 {
					m.Set(x, y, color.RGBA{255, 255, 255, 255})
				} else if fgBorder.RGBAAt(x, y).R == 0 {
					m.Set(x, y, color.RGBA{0, 0, 0, 255})
				}
			}
		}
	}
	return m
}

func AddWallpaperToFrame(fgImg *image.RGBA) *image.RGBA {
	if !config.Config.UseWallpaper {
		return fgImg
	}
	b := fgImg.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), BackgroundImage.(*image.RGBA), b.Min, draw.Src)
	return CompositeBackgroundAndForeground(m, fgImg)
}

func InvertImage(img *image.RGBA) *image.RGBA {
	return NRGBAImgToRGBAImg(imaging.Invert(img))
}

func MultiLinePixFont(img *image.RGBA, x, y, lineSep int, str string, col color.Color) {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		pixfont.DrawString(img, x, y+i*lineSep, line, col)
	}
}
