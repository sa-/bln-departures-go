package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	h "github.com/sa-/schedule/hafasClient"
)

type model struct {
	// source of truth
	departures   *h.DepartureBoard
	windowWidth  int
	windowHeight int

	// computed
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "r":
			table, departures := getData()
			return model{departures, m.windowWidth, m.windowWidth, table}, cmd
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

var tableStyle = lipgloss.NewStyle()

func (m model) View() string {
	windowStyle := lipgloss.NewStyle().Width(m.windowWidth).Height(m.windowHeight).Align(lipgloss.Center)
	renderedTable := tableStyle.Render(m.table.View())

	return windowStyle.Render(renderedTable)
}
