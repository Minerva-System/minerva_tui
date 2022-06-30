package hostscreen

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	api "luksamuk/minerva_tui/client"
	"net/http"
)

type Model struct {
	Client   api.MinervaClient
	Ready    bool
	status   string
	taborder int
	host     textinput.Model
	tenant   textinput.Model
	login    textinput.Model
	pass     textinput.Model
}

func Create() Model {
	// Hostname
	host := textinput.New()
	host.Placeholder = "http://localhost:9000"
	host.Focus()
	host.CharLimit = 200
	host.Width = 50
	host.Prompt = "Host:      "

	// Tenant
	tenant := textinput.New()
	tenant.Placeholder = "minerva"
	tenant.Blur()
	tenant.CharLimit = 50
	tenant.Width = 50
	tenant.Prompt = "Inquilino: "

	// Login
	login := textinput.New()
	login.Blur()
	login.CharLimit = 50
	login.Width = 50
	login.Prompt = "Login:     "

	// Password
	pass := textinput.New()
	pass.Blur()
	pass.CharLimit = 50
	pass.Width = 50
	pass.EchoMode = textinput.EchoPassword
	pass.Prompt = "Senha:     "

	return Model{
		Client:   api.MinervaClient{},
		Ready:    false,
		status:   "",
		taborder: 0,
		host:     host,
		tenant:   tenant,
		login:    login,
		pass:     pass,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) manageTabOrder(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	m.host.Blur()
	m.tenant.Blur()
	m.login.Blur()
	m.pass.Blur()

	switch m.taborder {
	case 0:
		m.host.Focus()
		m.host, cmd = m.host.Update(msg)
	case 1:
		m.tenant.Focus()
		m.tenant, cmd = m.tenant.Update(msg)
	case 2:
		m.login.Focus()
		m.login, cmd = m.login.Update(msg)
	case 3:
		m.pass.Focus()
		m.pass, cmd = m.pass.Update(msg)
	case 4:
	}

	return cmd
}

func (m *Model) reset() {
	m.host.SetValue("")
	m.tenant.SetValue("")
	m.login.SetValue("")
	m.pass.SetValue("")

	m.host.Focus()
	m.tenant.Blur()
	m.login.Blur()
	m.pass.Blur()
	
	m.taborder = 0
	m.Ready = false
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.taborder {
			case 0:
				m.host.SetValue(m.host.Placeholder)
			case 1:
				m.tenant.SetValue(m.tenant.Placeholder)
			}
		case "enter":
			switch m.taborder {
			case 0:
				if m.host.Value() == "" {
					m.host.SetValue(m.host.Placeholder)
				}
				m.status = ""
				m.taborder++
			case 1:
				if m.tenant.Value() == "" {
					m.tenant.SetValue(m.tenant.Placeholder)
				}
				m.taborder++
			case 2:
				if m.login.Value() == "" {
					m.status = "Informe o login do usuário."
				} else {
					m.status = ""
					m.taborder++
				}
			case 3:
				if m.pass.Value() == "" {
					m.status = "Informe a senha do usuário."
				} else {
					m.status = "Fazendo Login..."
					m.taborder++
				}
			}
		}
	}

	cmd = m.manageTabOrder(msg)

	if m.taborder == 4 {
		client, err := api.Create(m.host.Value(), m.tenant.Value())
		if err != nil {
			m.status = fmt.Sprintf("Erro: %v", err)
			return m, cmd
		}
		m.Client = client
		
		code, _, errmsg := m.Client.Login(api.LoginRequest{
			Login:    m.login.Value(),
			Password: m.pass.Value(),
		})

		if errmsg != "" {
			m.status = errmsg
		} else {
			if code != 200 {
				m.status = fmt.Sprintf("Erro: %d %s", code, http.StatusText(code))
			}
		}

		if code != 200 {
			m.reset()
		} else {
			m.Ready = true
		}
	}

	return m, cmd
}

func (m Model) View() string {
	s := "=== Minerva System ===\n\n"
	s += m.host.View() + "\n"
	s += m.tenant.View() + "\n"
	s += m.login.View() + "\n"
	s += m.pass.View() + "\n"
	s += "\n\n" + m.status + "\n\n"
	return s
}
