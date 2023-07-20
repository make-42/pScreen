package constants

const (
	ClockModuleID = iota
	BlankModuleID
	MediaModuleID
	MonitorModuleID
	PongModuleID
	ScreensaverModuleID
	WeatherModuleID
)

var AllModules = []int{
	BlankModuleID, ClockModuleID, MediaModuleID, MonitorModuleID, PongModuleID, ScreensaverModuleID, WeatherModuleID,
}

const (
	MainPageID = iota
	AddModulePageID
	BridgeSettingsPageID
)
