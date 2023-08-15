package constants

import "github.com/vishalkuo/bimap"

const AppVersion = "α 0.0.1"

const (
	ClockModuleID = iota
	BlankModuleID
	DiscordModuleID
	MediaModuleID
	MonitorModuleID
	OsuModuleID
	PongModuleID
	QRCodeModuleID
	ScreensaverModuleID
	WeatherModuleID
)

var AllModules = []int{
	BlankModuleID, ClockModuleID, DiscordModuleID, MediaModuleID, MonitorModuleID, OsuModuleID, PongModuleID, QRCodeModuleID, ScreensaverModuleID, WeatherModuleID,
}

var ModuleNames = bimap.NewBiMapFromMap[string, int](map[string]int{
	"blank":       BlankModuleID,
	"clock":       ClockModuleID,
	"discord":     DiscordModuleID,
	"media":       MediaModuleID,
	"monitor":     MonitorModuleID,
	"osu":         OsuModuleID,
	"pong":        PongModuleID,
	"qrcode":      QRCodeModuleID,
	"screensaver": ScreensaverModuleID,
	"weather":     WeatherModuleID,
})

const (
	MainPageID = iota
	AddModulePageID
	BridgeSettingsPageID
)
