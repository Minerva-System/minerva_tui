package main

import (
	"fmt"
	"luksamuk/minerva_tui/hostscreen"
	"luksamuk/minerva_tui/mainmenu"
	"luksamuk/minerva_tui/user_list"
	"luksamuk/minerva_tui/userform"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	ready       bool
	currentView int
	host        hostscreen.Model
	mainmenu    mainmenu.Model
	userlist    user_list.Model
	userform    userform.Model
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
		m.userlist.SetSize(msg.Width, msg.Height)
		m.userform.SetSize(msg.Width, msg.Height)
		m.ready = true
	}
	
	if m.ready {
		switch m.currentView {
		case 0:
			m.host, cmd = m.host.Update(msg)
			if m.host.Ready {
				m.currentView = 1
				m.mainmenu.Client = &m.host.Client
				m.userlist.Client = &m.host.Client
				m.userform.Client = &m.host.Client
				cmd = tea.Batch(cmd, tea.EnterAltScreen)
			}
		case 1:
			m.mainmenu, cmd = m.mainmenu.Update(msg)
			if m.mainmenu.Option == "Usuários" {
				m.currentView = 2
				m.mainmenu.Option = ""
				m.userlist.Fetch()
			}
		case 2:
			m.userlist, cmd = m.userlist.Update(msg)
			if m.userlist.Option != "" {
				switch m.userlist.Option {
				case "quit":
					m.currentView = 1
				case "create":
					m.currentView = 3
				}
				m.userlist.Option = ""
			}
		case 3:
			m.userform, cmd = m.userform.Update(msg)
			if m.userform.Option == "quit" {
				m.currentView = 2
				m.userform.Option = ""
				m.userlist.Fetch()
			}
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
		case 2:
			return m.userlist.View()
		case 3:
			return m.userform.View()
		}
	}
	return ""
}

func CreateApp() App {
	return App{
		currentView: 0,
		host:        hostscreen.Create(),
		mainmenu:    mainmenu.Create(),
		userlist:    user_list.Create(),
		userform:    userform.Create(),
	}
}

func main() {
	if err := tea.NewProgram(CreateApp(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Erro ao executar o programa:\n%v\n", err)
		os.Exit(1)
	}
}
