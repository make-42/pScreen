package qrcode

import (
	"image"
	"image/draw"
	"pscreenapp/bridge/modules"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"pscreenapp/utils"

	"git.sr.ht/~sbinet/gg"
	"github.com/skip2/go-qrcode"
)

var CurrentGIFFrame = 0

var QRCodeModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	qrcodeObj, err := qrcode.New(config.Config.QRCodeContent, qrcode.Highest)
	utils.CheckError(err)
	qrcodeImg := qrcodeObj.Image(-1)
	qrcodeImgBounds := qrcodeImg.Bounds()
	offset := image.Pt(4, config.Config.CanvasRenderDimensions.Y/2-qrcodeImgBounds.Dy()/2)
	draw.Draw(im, qrcodeImgBounds.Add(offset), qrcodeImg, image.Point{0, 0}, draw.Src)
	qrcodeGifBounds := renderer.QRCodeModuleGIF.Image[CurrentGIFFrame].Bounds()
	qrcodeGifOffset := image.Pt(qrcodeImgBounds.Dx()+8, config.Config.CanvasRenderDimensions.Y-qrcodeGifBounds.Dy())
	draw.Draw(im, qrcodeGifBounds.Add(qrcodeGifOffset), renderer.QRCodeModuleGIF.Image[CurrentGIFFrame], image.Point{0, 0}, draw.Src)
	dc.SetFontFace(renderer.MedSmallerFont)
	dc.DrawStringWrapped(config.Config.QRCodeTitle, float64(qrcodeImgBounds.Dy())+8, 4, 0, 0, float64(config.Config.CanvasRenderDimensions.X)-float64(qrcodeImgBounds.Dy())-12, 1, gg.AlignRight)
	dc.SetFontFace(renderer.TinyFont)
	dc.DrawStringWrapped(config.Config.QRCodeDescription, float64(qrcodeImgBounds.Dy())+8, 32, 0, 0, float64(config.Config.CanvasRenderDimensions.X)-float64(qrcodeImgBounds.Dy())-12, 1, gg.AlignRight)
	CurrentGIFFrame++
	if CurrentGIFFrame >= len(renderer.QRCodeModuleGIF.Image) {
		CurrentGIFFrame = 0
	}
	return renderer.RemoveAntiAliasing(im)
}}
