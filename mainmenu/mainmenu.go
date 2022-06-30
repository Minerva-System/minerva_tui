package mainmenu

import (
	"fmt"
	"net/http"
	"github.com/charmbracelet/bubbles/list"
	api "luksamuk/minerva_tui/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	Client  *api.MinervaClient
	list    list.Model
}

func Create() Model {
	items := []list.Item{
		item{title: "Usuários", desc: "Gerenciar usuários do sistema"},
		item{title: "Produtos", desc: "Gerenciar produtos"},
		item{title: "Estoque", desc: "Controles de estoque"},
		item{title: "Clientes", desc: "Gerenciar clientes"},
		item{title: "Auditoria", desc: "Pesquisar nos logs do sistema"},
		item{title: "CRM", desc: "Interação com o cliente"},
		item{title: "Relatórios", desc: "Visualizar relatórios"},
		item{title: "Logout", desc: "Sair do sistema"},
	}
	itemlist := list.New(items, list.NewDefaultDelegate(), 0, 0)
	itemlist.Title = "Menu Principal"
	itemlist.KeyMap.Quit.SetEnabled(false)
	itemlist.SetShowStatusBar(true)
	
	return Model{
		Client: nil,
		list: itemlist,
	}
}

func (m *Model) SetSize(width int, height int) {
	h, v := style.GetFrameSize()
	m.list.SetSize(width-h, height-v)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(item)			
			if ok {
				m.list.NewStatusMessage(i.Title())
				if i.Title() == "Logout" {
					code, _, errmsg := m.Client.Logout()
					if errmsg != "" {
						m.list.NewStatusMessage(errmsg)
					} else if code != 200 {
						m.list.NewStatusMessage(fmt.Sprintf("Erro ao efetuar logout: %s", http.StatusText(code)))
					} else {
						cmd = tea.Quit
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	var listcmd tea.Cmd
	m.list, listcmd = m.list.Update(msg)

	return m, tea.Batch(listcmd, cmd)
}

func (m Model) View() string {
	s := style.Render(m.list.View())
	return s
}
