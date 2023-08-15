package main

import (
	"log"
	"pscreen/bridge"
	"pscreen/bridge/comms"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/i18n"
	"pscreen/ui"
	"pscreen/utils"

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
