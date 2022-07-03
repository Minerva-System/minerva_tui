package userform

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	api "luksamuk/minerva_tui/client"
	"net/http"
	"strings"
)

type keyMap struct {
	Back    key.Binding
	Next    key.Binding
	Confirm key.Binding
	Mock    key.Binding
	Cancel  key.Binding
	Help    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Mock, k.Cancel, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Next, k.Confirm},
		{k.Mock, k.Cancel, k.Help},
	}
}

var (
	keys = keyMap{
		Back: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "anterior"),
		),
		Next: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "próximo"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "selecionar"),
		),
		Mock: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "preencher"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancelar"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "ajuda"),
		),
	}

	// Styles
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Model struct {
	Option   string
	Client   *api.MinervaClient
	keys     keyMap
	help     help.Model
	taborder int
	status   string
	login    textinput.Model
	name     textinput.Model
	email    textinput.Model
	pass     textinput.Model
}

func Create() Model {
	// Login
	login := textinput.New()
	login.Focus()
	login.CharLimit = 15
	login.Width = 15
	login.Prompt = "Login:   "

	// Name
	name := textinput.New()
	name.Blur()
	name.CharLimit = 80
	name.Width = 80
	name.Prompt = "Nome:    "

	// E-mail
	email := textinput.New()
	email.Blur()
	email.CharLimit = 50
	email.Width = 50
	email.Prompt = "E-mail:  "

	// Password
	pass := textinput.New()
	pass.Blur()
	pass.CharLimit = 50
	pass.Width = 50
	pass.EchoMode = textinput.EchoPassword
	pass.Prompt = "Senha:   "

	return Model{
		Client:   nil,
		Option:   "",
		keys:     keys,
		help:     help.New(),
		taborder: 0,
		status:   "",
		login:    login,
		name:     name,
		email:    email,
		pass:     pass,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetSize(width int, height int) {
	m.help.Width = width
}

func (m *Model) Reset() {
	m.taborder = 0
	m.login.SetValue("")
	m.name.SetValue("")
	m.email.SetValue("")
	m.pass.SetValue("")
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Cancel):
			m.Option = "quit"
		case key.Matches(msg, m.keys.Back):
			if m.taborder > 0 {
				m.taborder--
			}
		case key.Matches(msg, m.keys.Confirm):
			if m.taborder < 4 {
				if m.taborder == 0 {
					m.login.SetValue(
						strings.Replace(strings.ToLower(m.login.Value()), " ", "", -1))
				}
				m.taborder++
			} else if m.taborder == 4 {
				code, _, msg := m.Client.UserCreate(api.NewUserRequest{
					Login:    m.login.Value(),
					Name:     m.name.Value(),
					Email:    m.email.Value(),
					Password: m.pass.Value(),
				})

				if msg != "" {
					m.status = msg
				} else if code > 299 {
					m.status = fmt.Sprintf("%d %s", code, http.StatusText(code))
				} else {
					m.Reset()
					m.Option = "quit"
				}
			}
		case key.Matches(msg, m.keys.Next):
			if m.taborder < 4 {
				if m.taborder == 0 {
					m.login.SetValue(
						strings.Replace(strings.ToLower(m.login.Value()), " ", "", -1))
				}
				m.taborder++
			}
		case key.Matches(msg, m.keys.Mock):
			name := gofakeit.Name()
			username := strings.Replace(strings.ToLower(name), " ", "", -1)
			m.login.SetValue(username)
			m.name.SetValue(name)
			m.email.SetValue(gofakeit.Email())
			m.pass.SetValue("123456")
			m.taborder = 4
			m.status = "Dados preenchidos com sucesso."
		default:
			switch m.taborder {
			case 0:
				m.login, cmd = m.login.Update(msg)
			case 1:
				m.name, cmd = m.name.Update(msg)
			case 2:
				m.email, cmd = m.email.Update(msg)
			case 3:
				m.pass, cmd = m.pass.Update(msg)
			case 4:
			}
		}
	}

	// Switch focus
	m.login.Blur()
	m.name.Blur()
	m.email.Blur()
	m.pass.Blur()
	switch m.taborder {
	case 0:
		m.login.Focus()
	case 1:
		m.name.Focus()
	case 2:
		m.email.Focus()
	case 3:
		m.pass.Focus()
	}

	return m, cmd
}

func (m Model) View() string {
	s := m.login.View() + "\n\n"
	s += m.name.View() + "\n\n"
	s += m.email.View() + "\n\n"
	s += m.pass.View() + "\n\n\n"

	finalizar := "[ Finalizar Cadastro ]"
	if m.taborder == 4 {
		s += focusedStyle.Render(finalizar)
	} else {
		s += blurredStyle.Render(finalizar)
	}

	s += "\n\n\n"
	s += m.status + "\n"
	s += m.help.View(m.keys)
	return s
}
