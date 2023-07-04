package constants

const (
	ClockModuleID = iota
	MediaModuleID
	MonitorModuleID
	ScreensaverModuleID
	WeatherModuleID
)

var AllModules = []int{
	ClockModuleID, MediaModuleID, MonitorModuleID, ScreensaverModuleID, WeatherModuleID,
}

const (
	MainPageID = iota
	AddModulePageID
	BridgeSettingsPageID
)
