package constants

const (
	ClockModuleID = iota
	MediaModuleID
	ScreensaverModuleID
	WeatherModuleID
)

var AllModules = []int{
	ClockModuleID, MediaModuleID, ScreensaverModuleID, WeatherModuleID,
}

const (
	MainPageID = iota
	AddModulePageID
)
