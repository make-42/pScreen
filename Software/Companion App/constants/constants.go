package constants

const (
	ClockModuleID = iota
	BlankModuleID
	MediaModuleID
	MonitorModuleID
	ScreensaverModuleID
	WeatherModuleID
)

var AllModules = []int{
	BlankModuleID, ClockModuleID, MediaModuleID, MonitorModuleID, ScreensaverModuleID, WeatherModuleID,
}

const (
	MainPageID = iota
	AddModulePageID
	BridgeSettingsPageID
)
