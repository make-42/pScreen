package media

import (
	"image"
	"image/color"
	"image/draw"
	"net/url"
	"os"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/bridge/renderer/elements"
	"pscreen/config"
	"pscreen/utils"
	"runtime"
	"strings"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/disintegration/imaging"
	"github.com/godbus/dbus"
	"github.com/leberKleber/go-mpris"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"
)

var CurrentMediaArtURL = ""
var CurrentMediaArtImage *image.RGBA

var MediaModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	currentPlayingMediaInfo := GetCurrentPlayingMediaInfo()
	now := time.Now()
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.TinyFont)
	dc.DrawStringAnchored(currentPlayingMediaInfo.Title, 4, -2, 0, 1)
	dc.DrawStringAnchored(now.Format(currentPlayingMediaInfo.Album), 4, 10, 0, 1)
	dc.DrawStringAnchored(currentPlayingMediaInfo.Artist, 4, 22, 0, 1)
	elements.DrawMediaProgressBar(dc, currentPlayingMediaInfo.Position, currentPlayingMediaInfo.Duration)
	if CurrentMediaArtURL != "" {
		return renderer.CompositeBackgroundAndForeground(CurrentMediaArtImage, renderer.RemoveAntiAliasing(im))
	}
	return renderer.AddWallpaperToFrame(renderer.RemoveAntiAliasing(im))
}}

type CurrentPlayingMediaInfo struct {
	Title    string
	Album    string
	Artist   string
	Position float64
	Duration float64
}

func ListPlayers() []string {
	conn, err := dbus.SessionBus()
	utils.CheckError(err)
	var names []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	utils.CheckError(err)

	var mprisNames []string
	for _, name := range names {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2") {
			mprisNames = append(mprisNames, name)
		}
	}
	return mprisNames
}

func UpdateMediaArt(mediaArtURL string) {
	if mediaArtURL != CurrentMediaArtURL {
		CurrentMediaArtURL = mediaArtURL
		if CurrentMediaArtURL != "" {
			u, err := url.ParseRequestURI(CurrentMediaArtURL)
			utils.CheckError(err)
			imgFile, err := os.Open(u.Path)
			utils.CheckError(err)
			defer imgFile.Close()
			bgImg, _, err := image.Decode(imgFile)
			utils.CheckError(err)
			bgTImg := resize.Resize(uint(utils.Min(config.Config.CanvasRenderDimensions.X, config.Config.CanvasRenderDimensions.Y)), 0, bgImg, resize.Lanczos3)
			bgBImg := resize.Resize(uint(utils.Max(config.Config.CanvasRenderDimensions.X, config.Config.CanvasRenderDimensions.Y)), 0, bgImg, resize.Lanczos3)
			palette := []color.Color{
				color.Black,
				color.White,
			}
			d := dither.NewDitherer(palette)
			//d.Mapper = dither.Bayer(2, 2, 1.0)
			d.Matrix = dither.FloydSteinberg
			bgBImg = renderer.NRGBAImgToRGBAImg(imaging.Blur(bgBImg, 10))
			if !config.Config.RGBXMit {
				bgBImg = d.Dither(bgBImg)
				bgTImg = d.Dither(bgTImg)
			} else {
				bgTImg = renderer.YCbCrImgToRGBAImg(bgTImg.(*image.YCbCr))
			}

			tB := bgTImg.Bounds()
			bB := bgBImg.Bounds()
			m := image.NewRGBA(image.Rect(0, 0, config.Config.CanvasRenderDimensions.X, config.Config.CanvasRenderDimensions.Y))
			draw.Draw(m, bB.Bounds().Add(image.Pt((config.Config.CanvasRenderDimensions.X-bB.Dx())/2, (config.Config.CanvasRenderDimensions.Y-bB.Dy())/2)), bgBImg.(*image.RGBA), bB.Min, draw.Src)
			draw.Draw(m, tB.Bounds().Add(image.Pt((config.Config.CanvasRenderDimensions.X-tB.Dx()), (config.Config.CanvasRenderDimensions.Y-tB.Dy())/2)), bgTImg.(*image.RGBA), tB.Min, draw.Src)
			CurrentMediaArtImage = m
		}
	}
}

func GetCurrentPlayingMediaInfo() CurrentPlayingMediaInfo {
	switch runtime.GOOS {
	case "linux":
		players := ListPlayers()
		if len(players) == 0 {
			return CurrentPlayingMediaInfo{"No media", "", "", 0, 0}
		}
		p, err := mpris.NewPlayer(players[0])
		utils.CheckError(err)
		mediaPositionMicroseconds, err := p.Position()
		if err != nil {
			mediaPositionMicroseconds = 0
		}
		mediaPosition := float64(mediaPositionMicroseconds) / 1000000
		mediaMetadata, err := p.Metadata()
		utils.CheckError(err)
		mediaDurationMicroseconds, err := mediaMetadata.MPRISLength()
		utils.CheckError(err)
		mediaDuration := float64(mediaDurationMicroseconds) / 1000000
		mediaTitle, err := mediaMetadata.XESAMTitle()
		utils.CheckError(err)
		mediaAlbum, err := mediaMetadata.XESAMAlbum()
		utils.CheckError(err)
		mediaArtists, err := mediaMetadata.XESAMArtist()
		utils.CheckError(err)
		mediaArtist := ""
		if len(mediaArtists) != 0 {
			mediaArtist = mediaArtists[0]
		}
		mediaArtURL, err := mediaMetadata.MPRISArtURL()
		utils.CheckError(err)
		UpdateMediaArt(mediaArtURL)
		return CurrentPlayingMediaInfo{mediaTitle, mediaAlbum, mediaArtist, mediaPosition, mediaDuration}
	default:
		return CurrentPlayingMediaInfo{"Platform not supported", "Sorry", "", 0, 0}
	}
}
