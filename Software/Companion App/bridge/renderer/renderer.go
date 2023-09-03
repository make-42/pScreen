package renderer

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"os"
	"pscreen/bridge/encoder"
	"pscreen/bridge/modules"
	"pscreen/config"
	"pscreen/utils"
	"strings"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/creasty/go-easing"
	"github.com/disintegration/imaging"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/pbnjay/pixfont"
	"golang.org/x/image/font"
)

var SavedFrames int

//go:embed assets/fonts/iosevka-heavy.ttf
//go:embed assets/fonts/iosevka-medium.ttf
//go:embed assets/img/bg.png
//go:embed assets/img/DVD_logo.png
//go:embed assets/gif/komi-san-48.gif
var assetsFolder embed.FS

var BigFont font.Face
var MediumFont font.Face
var MedSmallFont font.Face
var MedSmallerFont font.Face
var SmallFont font.Face
var TinyFont font.Face
var BackgroundImage image.Image
var DVDLogo image.Image
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
	utils.CheckError(err)
	defer imgFile.Close()
	BackgroundImage, _, err = image.Decode(imgFile)
	utils.CheckError(err)
	dvdImgFile, err := assetsFolder.Open("assets/img/DVD_logo.png")
	utils.CheckError(err)
	defer dvdImgFile.Close()
	DVDLogo, _, err = image.Decode(dvdImgFile)
	utils.CheckError(err)
	gifFile, err := assetsFolder.Open("assets/gif/komi-san-48.gif")
	utils.CheckError(err)
	defer gifFile.Close()
	QRCodeModuleGIF, err = gif.DecodeAll(gifFile)
	utils.CheckError(err)
}

func RenderFrame(mod modules.Module, lastMod modules.Module, timeOfSwitch int64) []byte {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{config.Config.CanvasRenderDimensions.X, config.Config.CanvasRenderDimensions.Y}
	im := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	transitionAmount := math.Min(float64(time.Now().UTC().UnixMilli()-timeOfSwitch)/config.Config.TransitionMilliseconds, 1)
	if transitionAmount < 0.5 {
		im = lastMod.RenderFunction(im)
		im = BlurImage(im, easing.Transition(0, config.Config.TransitionBlurringSigma, easing.QuadEaseInOut(transitionAmount*2)))
	} else {
		im = mod.RenderFunction(im)
		im = BlurImage(im, easing.Transition(config.Config.TransitionBlurringSigma, 0, easing.QuadEaseInOut((transitionAmount-0.5)*2)))
	}

	if config.Config.DebugSaveScreen {
		f, _ := os.Create(fmt.Sprintf("frame-%04d.png", SavedFrames))
		png.Encode(f, im)
		SavedFrames++
	}
	if config.Config.RotateScreen180Deg {
		im = RotateImage180deg(im)
	}
	return encoder.EncodeFrameToBytes(im)

}

func BlurImage(im *image.RGBA, sigma float64) *image.RGBA {
	if sigma == 0 {
		return im
	}
	im = NRGBAImgToRGBAImg(imaging.Blur(im, sigma))
	if !config.Config.RGBXMit {
		palette := []color.Color{
			color.Black,
			color.White,
		}
		d := dither.NewDitherer(palette)
		//d.Mapper = dither.Bayer(2, 2, 1.0)
		d.Matrix = dither.FloydSteinberg
		return d.Dither(im).(*image.RGBA)
	} else {
		return im
	}
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
