package ui

import (
	"fmt"
	"pscreenapp/bridge"
	"pscreenapp/config"
	"pscreenapp/constants"
	"pscreenapp/i18n"
	"pscreenapp/ui/components"
	"pscreenapp/ui/styling"
	"pscreenapp/utils"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	currentPage int
	cursor      int
}

func InitialModel() Model {
	return Model{
		currentPage: constants.MainPageID,
		cursor:      0,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		var navItems int
		switch m.currentPage {
		case constants.MainPageID:
			navItems = len(bridge.BridgeData.LoadedModules)
		case constants.AddModulePageID:
			navItems = len(constants.AllModules)
		case constants.BridgeSettingsPageID:
			navItems = len(bridge.BridgeData.DetectedPorts)
		}
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < navItems-1 {
				m.cursor++
			}

		// The "a" key adds a new module
		case "a":
			if m.currentPage == constants.MainPageID {
				m.currentPage = constants.AddModulePageID
				m.cursor = 0
			}

		// The "r" key removes a module
		case "r":
			if m.currentPage == constants.MainPageID {
				if len(bridge.BridgeData.LoadedModules) > 0 {
					bridge.BridgeData.LoadedModules = utils.DeleteItem(bridge.BridgeData.LoadedModules, m.cursor)
					if m.cursor > 0 {
						m.cursor--
					}
				}
			}

		// The "s" key starts communication with the device
		case "s":
			if m.currentPage == constants.MainPageID {
				if !bridge.BridgeData.CommsReady {
					bridge.BridgeStartXMit()
				}
			}

		// The "b" opens bridge settings page
		case "b":
			if m.currentPage == constants.MainPageID {
				m.cursor = 0
				m.currentPage = constants.BridgeSettingsPageID
				bridge.BridgeEnumSerialDevices()
			}

		// The "enter" and "space" keys select an item
		case "enter", " ":
			switch m.currentPage {
			case constants.AddModulePageID:
				bridge.BridgeData.LoadedModules = append(bridge.BridgeData.LoadedModules, constants.AllModules[m.cursor])
				m.currentPage = constants.MainPageID
				m.cursor = 0
			case constants.BridgeSettingsPageID:
				if len(bridge.BridgeData.DetectedPorts) > 0 {
					config.SerialPortToUse = bridge.BridgeData.DetectedPorts[m.cursor].Name
					m.currentPage = constants.MainPageID
					m.cursor = 0
				}
			}

		// The "esc" key leave a menu
		case "esc":
			if m.currentPage != constants.MainPageID {
				m.currentPage = constants.MainPageID
				m.cursor = 0
			}
		}
	case tickMsg:
		return m, tickCmd()
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	s := "\n"
	s += components.Header(components.VersionNumber())
	s += "\n\n"
	keybinds := []string{}
	switch m.currentPage {
	case constants.MainPageID:
		// The header
		s += components.Header(i18n.I18n.Headers.MainPageHeader) + "\n\n"

		// Iterate over our modules
		for i, moduleID := range bridge.BridgeData.LoadedModules {

			// Is the cursor pointing at this module?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = styling.ColorFg(">", styling.HighlightedColor) // cursor!
			}

			shownOnScreen := " "
			if bridge.BridgeData.CurrentModule == i {
				shownOnScreen = styling.ColorFg("*", styling.IndicatorColor)
			}

			// Render the row
			s += styling.Indent(fmt.Sprintf("%s %s %s\n", shownOnScreen, cursor, i18n.ModuleIDToString(moduleID)), config.PaddingIndentAmount)
		}
		keybinds = append(keybinds, i18n.I18n.Keybinds.AddModuleKeybind, i18n.I18n.Keybinds.RemoveModuleKeybind, i18n.I18n.Keybinds.StartXMitKeybind, i18n.I18n.Keybinds.BridgeSettingsKeybind)

	case constants.AddModulePageID:
		// The header
		s += components.Header(i18n.I18n.Headers.AddModulePageHeader) + "\n\n"

		// Iterate over our modules
		for i, moduleID := range constants.AllModules {

			// Is the cursor pointing at this module?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = styling.ColorFg(">", styling.HighlightedColor) // cursor!
			}

			// Render the row
			s += styling.Indent(fmt.Sprintf("%s %s\n", cursor, i18n.ModuleIDToString(moduleID)), config.PaddingIndentAmount)
		}
		keybinds = append(keybinds, i18n.I18n.Keybinds.SelectKeybind, i18n.I18n.Keybinds.EscKeybind)

	case constants.BridgeSettingsPageID:
		// The header
		s += components.Header(i18n.I18n.Headers.BridgeSettingsPageHeader) + "\n\n"

		// Iterate over our modules
		for i, port := range bridge.BridgeData.DetectedPorts {

			// Is the cursor pointing at this module?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = styling.ColorFg(">", styling.HighlightedColor) // cursor!
			}
			portName := port.Name
			if config.SerialPortToUse == portName {
				portName = styling.Bold(portName)
			}
			// Render the row
			s += styling.Indent(fmt.Sprintf("%s %-10s %10s %10s \n", cursor, portName, port.VID, port.SerialNumber), config.PaddingIndentAmount)
		}
		keybinds = append(keybinds, i18n.I18n.Keybinds.SelectKeybind, i18n.I18n.Keybinds.EscKeybind)

	}
	keybinds = append(keybinds, i18n.I18n.Keybinds.NavigateKeybind, i18n.I18n.Keybinds.ExitKeybind)
	// The footer
	s += "\n" + components.KeybindsHints(keybinds) + "\n"

	// Send the UI for rendering
	return s
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*config.UpdateUIEveryXMilliseconds/1000, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
