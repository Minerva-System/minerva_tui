// Default package for `minerva_tui`.
package main

import (
	"fmt"
	table "github.com/calyptia/go-bubble-table"
	"luksamuk/minerva_tui/hostscreen"
	"luksamuk/minerva_tui/mainmenu"
	"luksamuk/minerva_tui/userform"
	"luksamuk/minerva_tui/userlist"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// Model for the application.
type App struct {
	ready       bool
	currentView int
	host        hostscreen.Model
	mainmenu    mainmenu.Model
	userlist    userlist.Model
	userform    userform.Model
}

// Function for Bubble Tea initialization of the model.
func (m App) Init() tea.Cmd {
	return nil
}

// Function for updating the screens on events.
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
				case "edit":
					row := m.userlist.Table.SelectedRow().(table.SimpleRow)
					index, _ := strconv.ParseInt(row[0].(string), 10, 64)
					m.userform.PrepareEdit(
						index,
						row[1].(string), // login
						row[2].(string), // name
						row[3].(string), // email
					)
					m.currentView = 3
					m.userlist.Option = ""
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

// Function for rendering the screens.
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

// Creates a new application model.
func CreateApp() App {
	return App{
		currentView: 0,
		host:        hostscreen.Create(),
		mainmenu:    mainmenu.Create(),
		userlist:    userlist.Create(),
		userform:    userform.Create(),
	}
}

// Main entry point.
func main() {
	if err := tea.NewProgram(CreateApp(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Erro ao executar o programa:\n%v\n", err)
		os.Exit(1)
	}
}
