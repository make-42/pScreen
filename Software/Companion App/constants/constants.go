package constants

const (
	ClockModuleID = iota
	BlankModuleID
	MediaModuleID
	MonitorModuleID
	PongModuleID
	QRCodeModuleID
	ScreensaverModuleID
	WeatherModuleID
)

var AllModules = []int{
	BlankModuleID, ClockModuleID, MediaModuleID, MonitorModuleID, PongModuleID, QRCodeModuleID, ScreensaverModuleID, WeatherModuleID,
}

const (
	MainPageID = iota
	AddModulePageID
	BridgeSettingsPageID
)
