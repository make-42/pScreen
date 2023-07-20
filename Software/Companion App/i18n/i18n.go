package i18n

import "pscreenapp/constants"

type ModuleStrings struct {
	BlankModule       string
	ClockModule       string
	MediaModule       string
	MonitorModule     string
	PongModule        string
	ScreensaverModule string
	WeatherModule     string
}

func ModuleIDToString(moduleID int) string {
	switch moduleID {
	case constants.BlankModuleID:
		return I18n.Modules.BlankModule
	case constants.ClockModuleID:
		return I18n.Modules.ClockModule
	case constants.MediaModuleID:
		return I18n.Modules.MediaModule
	case constants.MonitorModuleID:
		return I18n.Modules.MonitorModule
	case constants.PongModuleID:
		return I18n.Modules.PongModule
	case constants.ScreensaverModuleID:
		return I18n.Modules.ScreensaverModule
	case constants.WeatherModuleID:
		return I18n.Modules.WeatherModule
	default:
		return I18n.CatchAll
	}
}

type HeaderStrings struct {
	MainPageHeader           string
	AddModulePageHeader      string
	BridgeSettingsPageHeader string
}

type KeybindStrings struct {
	NavigateKeybind       string
	AddModuleKeybind      string
	RemoveModuleKeybind   string
	StartXMitKeybind      string
	SelectKeybind         string
	BridgeSettingsKeybind string
	EscKeybind            string
	ExitKeybind           string
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
	case "fr_FR":
		I18n = fr_FR
	default:
		I18n = en_US
	}
}

var en_US = LanguageStrings{
	Headers: HeaderStrings{
		MainPageHeader:           "Configuration",
		AddModulePageHeader:      "Add module",
		BridgeSettingsPageHeader: "Bridge settings",
	},
	Modules: ModuleStrings{BlankModule: "Blank",
		ClockModule:       "Clock",
		MediaModule:       "Media",
		MonitorModule:     "Monitor",
		PongModule:        "Pong",
		ScreensaverModule: "Screensaver",
		WeatherModule:     "Weather",
	},
	Keybinds: KeybindStrings{
		NavigateKeybind:       "j/k, up/down: select",
		AddModuleKeybind:      "a: add module",
		RemoveModuleKeybind:   "r: remove module",
		StartXMitKeybind:      "s: start xmit",
		SelectKeybind:         "enter/space: select",
		BridgeSettingsKeybind: "b: bridge settings",
		EscKeybind:            "esc: back",
		ExitKeybind:           "q: quit",
	},
	CatchAll: "Err;i18n",
}

var fr_FR = LanguageStrings{
	Headers: HeaderStrings{
		MainPageHeader:           "Configuration",
		AddModulePageHeader:      "Ajouter un module",
		BridgeSettingsPageHeader: "Réglages du pont",
	},
	Modules: ModuleStrings{
		BlankModule:       "Vide",
		ClockModule:       "Horloge",
		MediaModule:       "Média",
		MonitorModule:     "Moniteur",
		PongModule:        "Pong",
		ScreensaverModule: "Économiseur d'écran",
		WeatherModule:     "Météo",
	},
	Keybinds: KeybindStrings{
		NavigateKeybind:       "j/k, haut/bas: sélectionner",
		AddModuleKeybind:      "a: ajouter un module",
		RemoveModuleKeybind:   "r: supprimer le module",
		StartXMitKeybind:      "s: démarrer xmition",
		SelectKeybind:         "entrée/espace: sélectionner",
		BridgeSettingsKeybind: "b: réglages du pont",
		EscKeybind:            "esc: arrière",
		ExitKeybind:           "q: quitter",
	},
	CatchAll: "Err;i18n",
}
