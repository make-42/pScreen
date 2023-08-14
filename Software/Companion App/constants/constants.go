package constants

const (
	ClockModuleID = iota
	BlankModuleID
	DiscordModuleID
	MediaModuleID
	MonitorModuleID
	PongModuleID
	QRCodeModuleID
	ScreensaverModuleID
	WeatherModuleID
)

var AllModules = []int{
	BlankModuleID, ClockModuleID, DiscordModuleID, MediaModuleID, MonitorModuleID, PongModuleID, QRCodeModuleID, ScreensaverModuleID, WeatherModuleID,
}

const (
	MainPageID = iota
	AddModulePageID
	BridgeSettingsPageID
)
