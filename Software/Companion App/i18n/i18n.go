package i18n

import "pscreenapp/constants"

type ModuleStrings struct {
	ClockModule       string
	MediaModule       string
	ScreensaverModule string
	WeatherModule     string
}

func ModuleIDToString(moduleID int) string {
	switch moduleID {
	case constants.ClockModuleID:
		return I18n.Modules.ClockModule
	case constants.MediaModuleID:
		return I18n.Modules.MediaModule
	case constants.ScreensaverModuleID:
		return I18n.Modules.ScreensaverModule
	case constants.WeatherModuleID:
		return I18n.Modules.WeatherModule
	default:
		return I18n.CatchAll
	}
}

type HeaderStrings struct {
	MainPageHeader      string
	AddModulePageHeader string
}

type KeybindStrings struct {
	NavigateKeybind     string
	AddModuleKeybind    string
	RemoveModuleKeybind string
	SelectKeybind       string
	EscKeybind          string
	ExitKeybind         string
}

type LanguageStrings struct {
	Headers  HeaderStrings
	Modules  ModuleStrings
	Keybinds KeybindStrings
	CatchAll string
}

var I18n LanguageStrings

func LoadLanguageStrings(languageString string) {
	switch languageString {
	case "en_US":
		I18n = en_US
	default:
		I18n = en_US
	}
}

var en_US = LanguageStrings{
	Headers: HeaderStrings{
		MainPageHeader:      "Configuration",
		AddModulePageHeader: "Add module",
	},
	Modules: ModuleStrings{ClockModule: "Clock", MediaModule: "Media", ScreensaverModule: "Screensaver", WeatherModule: "Weather"},
	Keybinds: KeybindStrings{
		NavigateKeybind:     "j/k, up/down: select",
		AddModuleKeybind:    "a: add module",
		RemoveModuleKeybind: "r: remove module",
		SelectKeybind:       "enter/space: select",
		EscKeybind:          "esc: back",
		ExitKeybind:         "q: quit",
	},
	CatchAll: "Err;i18n",
}
