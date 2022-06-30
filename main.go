package main

import (
	"fmt"
	"os"
	"luksamuk/minerva_tui/hostscreen"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	currentView int
	host        hostscreen.Model
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
	}

	switch m.currentView {
	case 0:
		m.host, cmd = m.host.Update(msg)
	}

	return m, cmd
}

func (m App) View() string {
	switch m.currentView {
	case 0:
		return m.host.View()
	}

	return ""
}

func CreateApp() App {
	return App{
		currentView: 0,
		host: hostscreen.Create(),
	}
}

func main() {
	if err := tea.NewProgram(CreateApp()).Start(); err != nil {
		fmt.Printf("Erro ao executar o programa:\n%v\n", err)
		os.Exit(1)
	}
}
