package screen

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	SetSize(width int, height int)
	manageTabOrder(msg tea.Msg) tea.Cmd
	Reset()
	Init() tea.Cmd
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
}
