package main

import (
	"log"
	"pscreenapp/bridge"
	"pscreenapp/bridge/comms"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"pscreenapp/i18n"
	"pscreenapp/ui"
	"pscreenapp/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config.ParseConfig()
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	i18n.LoadLanguageStrings(config.Config.I18nLanguage)
	renderer.LoadRendererSharedRessources()
	comms.SetDefaultSerialPort()
	go bridge.BridgeMainThread()
	if config.Config.AutoStartXMit {
		bridge.BridgeStartXMit()
	}
	model := ui.InitialModel()
	p := tea.NewProgram(model)
	_, err := p.Run()
	utils.CheckError(err)
}
