package weather

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"net/http"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/utils"
	"time"

	"git.sr.ht/~sbinet/gg"

	"github.com/ztrue/tracerr"
)

var WeatherModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	if time.Now().UTC().UnixMilli()-LastTimeWeatherDataWasUpdated > config.Config.UpdateWeatherEveryXMilliseconds {
		RequestWeatherData()
		LastTimeWeatherDataWasUpdated = time.Now().UTC().UnixMilli()
	}
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.MediumFont)
	dc.DrawStringAnchored(fmt.Sprintf("%0.1f°C", KelvinToCelsius(CurrentWeatherData.Main.Temp)), 4, -6, 0, 1)
	dc.DrawStringAnchored(CurrentWeatherData.Weather[0].Main, float64(config.Config.CanvasRenderDimensions.X-4), -6, 1, 1)
	dc.SetFontFace(renderer.SmallFont)
	dc.DrawStringAnchored(fmt.Sprintf("%0.1f  %0.1f", KelvinToCelsius(CurrentWeatherData.Main.TempMin), KelvinToCelsius(CurrentWeatherData.Main.TempMax)), 4, 26, 0, 1)
	dc.DrawStringAnchored(fmt.Sprintf("%d hPa", CurrentWeatherData.Main.Pressure), float64(config.Config.CanvasRenderDimensions.X-4), 26, 1, 1)
	dc.DrawStringAnchored(fmt.Sprintf("%d %%", CurrentWeatherData.Main.Humidity), 4, float64(config.Config.CanvasRenderDimensions.Y-4), 0, 0)
	dc.DrawStringAnchored(fmt.Sprintf("%0.1f m/s", CurrentWeatherData.Wind.Speed), float64(config.Config.CanvasRenderDimensions.X)-4*2-config.Config.WindIndicatorRadius*2, float64(config.Config.CanvasRenderDimensions.Y-4), 1, 0)
	windDirectionLogoCenter := [2]float64{float64(config.Config.CanvasRenderDimensions.X) - 4 - config.Config.WindIndicatorRadius, float64(config.Config.CanvasRenderDimensions.Y) - 4 - config.Config.WindIndicatorRadius}
	radWindDir := float64(CurrentWeatherData.Wind.Deg) / 180 * math.Pi
	windDirectionWindEnd := [2]float64{windDirectionLogoCenter[0] + math.Sin(radWindDir)*config.Config.WindIndicatorRadius, windDirectionLogoCenter[1] - math.Cos(radWindDir)*config.Config.WindIndicatorRadius}
	dc.SetLineWidth(1)
	dc.DrawLine(windDirectionLogoCenter[0], windDirectionLogoCenter[1], windDirectionWindEnd[0], windDirectionWindEnd[1])
	dc.DrawCircle(windDirectionLogoCenter[0], windDirectionLogoCenter[1], config.Config.WindIndicatorRadius)
	dc.Stroke()
	return renderer.RemoveAntiAliasing(im)
}}

func KelvinToCelsius(K float64) float64 {
	return K - 273.15
}

var CurrentWeatherData WeatherData
var LastTimeWeatherDataWasUpdated int64

type WeatherDataWeather struct {
	ID          uint16
	Main        string
	Description string
	Icon        string
}

type MainWeatherData struct {
	Temp        float64
	FeelsLike   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Pressure    uint16
	Humidity    uint8
	SeaLevel    uint16 `json:"sea_level"`
	GroundLevel uint16 `json:"grnd_level"`
}

type WeatherCoords struct {
	Lat float64
	Lon float64
}

type WeatherWindData struct {
	Speed float64
	Deg   int16
	Gust  float64
}

type HistoricalPrecipitationData struct {
	OneH   float64 `json:"1h"`
	ThreeH float64 `json:"3h"`
}

type CloudsWeatherData struct {
	All uint8
}

type SysWeatherData struct {
	Type    uint64
	Id      uint64
	Country string
	Sunrise uint64
	Sunset  uint64
}

type WeatherData struct {
	Coord          WeatherCoords
	Weather        []WeatherDataWeather
	Base           string
	Timezone       int64
	TimezoneOffset int64
	Main           MainWeatherData
	Visibility     uint16
	Wind           WeatherWindData
	Rain           HistoricalPrecipitationData
	Snow           HistoricalPrecipitationData
	Clouds         CloudsWeatherData
	DT             uint64
	Sys            SysWeatherData
	Id             uint64
	Name           string
	Cod            uint64
}

func RequestWeatherData() {
	res, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", config.Config.Lat, config.Config.Long, config.Config.OpenWeatherMapAPIKey))
	utils.CheckError(tracerr.Wrap(err))
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&CurrentWeatherData)
	utils.CheckError(tracerr.Wrap(err))
}
