package constants

import "github.com/vishalkuo/bimap"

const AppVersion = "α 0.0.2"

const (
	ClockModuleID = iota
	BlankModuleID
	DiscordModuleID
	DVDModuleID
	MediaModuleID
	MonitorModuleID
	OsuModuleID
	PongModuleID
	QRCodeModuleID
	ScreensaverModuleID
	TeapotModuleID
	VisualizerModuleID
	WeatherModuleID
)

var AllModules = []int{
	BlankModuleID, ClockModuleID, DiscordModuleID, DVDModuleID, MediaModuleID, MonitorModuleID, OsuModuleID, PongModuleID, QRCodeModuleID, ScreensaverModuleID, TeapotModuleID, VisualizerModuleID, WeatherModuleID,
}

var ModuleNames = bimap.NewBiMapFromMap[string, int](map[string]int{
	"blank":       BlankModuleID,
	"clock":       ClockModuleID,
	"discord":     DiscordModuleID,
	"dvd":         DVDModuleID,
	"media":       MediaModuleID,
	"monitor":     MonitorModuleID,
	"osu":         OsuModuleID,
	"pong":        PongModuleID,
	"qrcode":      QRCodeModuleID,
	"screensaver": ScreensaverModuleID,
	"teapot":      TeapotModuleID,
	"visualizer":  VisualizerModuleID,
	"weather":     WeatherModuleID,
})

const (
	MainPageID = iota
	AddModulePageID
	BridgeSettingsPageID
)
