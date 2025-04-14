package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	h "github.com/sa-/schedule/hafasClient"
)

var appStateDepartureBoard *h.DepartureBoard

const displayModeCount = 2

type model struct {
	// source of truth
	windowWidth  int
	windowHeight int
	displayMode  int

	// computed
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyRight.String():
			m.displayMode = (m.displayMode + 1) % displayModeCount
		case tea.KeyLeft.String():
			m.displayMode = (m.displayMode - 1 + displayModeCount) % displayModeCount
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	windowStyle := lipgloss.NewStyle().Width(m.windowWidth).Height(m.windowHeight).Align(lipgloss.Center)
	renderedTable := m.table.View()
	switch m.displayMode {
	case 0:
		return windowStyle.Render(renderedTable)
	case 1:
		content := lipgloss.NewStyle().Align(lipgloss.Left).Render(simpleList(appStateDepartureBoard))
		return windowStyle.Render(content)
	}
	return windowStyle.Render("Invalid display mode: " + strconv.Itoa(m.displayMode))
}
