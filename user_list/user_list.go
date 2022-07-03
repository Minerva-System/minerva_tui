package user_list

import (
	"fmt"
	"strconv"
	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	api "luksamuk/minerva_tui/client"
)

type keyMap struct {
	Up           key.Binding
	Down         key.Binding
	Edit         key.Binding
	Create       key.Binding
	Remove       key.Binding
	NextPage     key.Binding
	PreviousPage key.Binding
	Help         key.Binding
	Quit         key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Create, k.Edit, k.Remove, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.NextPage, k.PreviousPage},
		{k.Create, k.Edit, k.Remove},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "navegar acima"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "navegar abaixo"),
	),
	Edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "editar"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "criar"),
	),
	Remove: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "remover"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "próxima página"),
	),
	PreviousPage: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "página anterior"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "ajuda"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "voltar"),
	),
}

type Model struct {
	Option    string
	Client    *api.MinervaClient
	keys      keyMap
	help      help.Model
	table     table.Model
	page      int
	maxPages  int
	paginator paginator.Model
	status    string
}

func Create() Model {
	table := table.New([]string{"ID", "LOGIN", "NOME", "E-MAIL"}, 10, 10)
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 1
	p.SetTotalPages(1)

	return Model{
		Option:    "",
		Client:    nil,
		keys:      keys,
		help:      help.New(),
		table:     table,
		page:      0,
		maxPages:  1,
		paginator: p,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Fetch() {
	code, data, errmsg := m.Client.UserList(m.page)
	if errmsg != "" {
		return
	}

	if code != 200 {
		return
	}

	if len(data) > 0 {
		rows := make([]table.Row, len(data))
		for i := 0; i < len(data); i++ {
			rows[i] = table.SimpleRow{
				fmt.Sprintf("%08d", data[i].ID),
				data[i].Login,
				data[i].Name,
				data[i].Email,
			}
		}
		m.table.SetRows(rows)
		m.table.GoPageUp()
		m.status = ""
	}

	if (len(data) > 0) && (m.maxPages < (m.page + 1)) {
		m.maxPages = m.page + 1
	} else if (len(data) == 0) {
		if m.page > 0 {
			m.page--
			m.Fetch()
		}
		m.maxPages = m.page + 1
	}
}

func (m *Model) SetSize(width int, height int) {
	diff := 3
	if m.help.ShowAll {
		diff = 6
	}

	m.table.SetSize(width, height-diff)
	m.help.Width = width
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.Option = "quit"
		case key.Matches(msg, m.keys.Create):
			m.Option = "create"
		case key.Matches(msg, m.keys.NextPage):
			m.page++
			m.Fetch()
		case key.Matches(msg, m.keys.PreviousPage):
			if m.page > 0 {
				m.page--
				m.Fetch()
			}
		case key.Matches(msg, m.keys.Remove):
			row := m.table.SelectedRow().(table.SimpleRow)
			index, err := strconv.ParseInt(row[0].(string), 10, 64)
			if row[1].(string) == "admin" {
				m.status = "Erro: Não é possível remover o administrador do sistema."
			} else if err != nil {
				m.status = fmt.Sprintf("Erro: %v", err)
			} else {
				_, errMsg := m.Client.UserRemove(index)
				if errMsg != "" {
					m.status = errMsg
				} else {
					m.Fetch()
					m.status = fmt.Sprintf("\"%s\" removido.", row[1].(string))
				}
			}
		}
	}

	m.paginator.SetTotalPages(m.maxPages)
	m.paginator.Page = m.page

	var tablecmd tea.Cmd
	m.table, tablecmd = m.table.Update(msg)
	return m, tea.Batch(cmd, tablecmd)
}

func (m Model) View() string {
	s := m.table.View() + "\n"
	s += m.paginator.View() + "\n"
	s += m.status + "\n"
	s += m.help.View(m.keys)
	return s
}
