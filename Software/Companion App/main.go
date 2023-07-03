package main

import (
	"pscreenapp/bridge"
	"pscreenapp/bridge/renderer"
	"pscreenapp/i18n"
	"pscreenapp/ui"
	"pscreenapp/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	i18n.LoadLanguageStrings("en_US")
	renderer.LoadRendererSharedRessources()
	go bridge.BridgeMainThread()
	model := ui.InitialModel()
	p := tea.NewProgram(model)
	_, err := p.Run()
	utils.CheckError(err)
}
