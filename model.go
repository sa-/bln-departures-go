package main

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
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
	departureTable table.Model
	hourlyViewport viewport.Model
	dailyViewport  viewport.Model

	altViewport viewport.Model
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
		m.hourlyViewport.Width = 20
		m.hourlyViewport.Height = msg.Height
		m.dailyViewport.Width = 50
		m.dailyViewport.Height = msg.Height
		m.altViewport.Width = msg.Width
		m.altViewport.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyTab.String():
			m.displayMode = (m.displayMode + 1) % displayModeCount
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.departureTable.SelectedRow()[1]),
			)
		}
	}
	if m.displayMode == 0 {
		m.departureTable, cmd = m.departureTable.Update(msg)
		m.hourlyViewport, cmd = m.hourlyViewport.Update(msg)
		m.dailyViewport, cmd = m.dailyViewport.Update(msg)
	} else if m.displayMode == 1 {
		m.altViewport, cmd = m.altViewport.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	windowStyle := lipgloss.NewStyle().Width(m.windowWidth).Height(m.windowHeight).Align(lipgloss.Center)

	switch m.displayMode {
	case 0:
		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.departureTable.View(),
			strings.Repeat(" | \n", m.windowHeight),
			m.hourlyViewport.View(),
			m.dailyViewport.View(),
		)
		return windowStyle.Render(content)
	case 1:
		return windowStyle.Render(m.altViewport.View())
	}
	return windowStyle.Render("Invalid display mode: " + strconv.Itoa(m.displayMode))
}
