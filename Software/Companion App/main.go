package main

import (
	"log"
	"pscreenapp/bridge"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"pscreenapp/i18n"
	"pscreenapp/ui"
	"pscreenapp/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	i18n.LoadLanguageStrings(config.I18nLanguage)
	renderer.LoadRendererSharedRessources()
	go bridge.BridgeMainThread()
	model := ui.InitialModel()
	p := tea.NewProgram(model)
	_, err := p.Run()
	utils.CheckError(err)
}
