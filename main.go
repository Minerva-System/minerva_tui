package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"luksamuk/minerva_tui/hostscreen"
	"luksamuk/minerva_tui/mainmenu"
	"os"
)

type App struct {
	ready       bool
	currentView int
	host        hostscreen.Model
	mainmenu    mainmenu.Model
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.mainmenu.SetSize(msg.Width, msg.Height)
		m.ready = true
	}
	
	if m.ready {
		switch m.currentView {
		case 0:
			m.host, cmd = m.host.Update(msg)
			if m.host.Ready {
				m.currentView++
				m.mainmenu.Client = &m.host.Client
				cmd = tea.Batch(cmd, tea.EnterAltScreen)
			}
		case 1:
			m.mainmenu, cmd = m.mainmenu.Update(msg)
		}
	}

	return m, cmd
}

func (m App) View() string {
	if m.ready {
		switch m.currentView {
		case 0:
			return m.host.View()
		case 1:
			return m.mainmenu.View()
		}
	}
	return ""
}

func CreateApp() App {
	return App{
		currentView: 0,
		host:        hostscreen.Create(),
		mainmenu:    mainmenu.Create(),
	}
}

func main() {
	if err := tea.NewProgram(CreateApp(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Erro ao executar o programa:\n%v\n", err)
		os.Exit(1)
	}
}
