package user_list

import (
	api "luksamuk/minerva_tui/client"
	tea "github.com/charmbracelet/bubbletea"
	table "github.com/calyptia/go-bubble-table"
)

type Model struct {
	Option string
	Client *api.MinervaClient
	table  table.Model
}

func Create() Model {
	table := table.New([]string{"ID", "LOGIN", "NOME", "E-MAIL"}, 10, 10)
	return Model{
		Option: "",
		Client: nil,
		table: table,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Fetch() {
	code, data, errmsg := m.Client.UserList(0)
	if errmsg != "" {
		return
	}

	if code != 200 {
		return
	}

	rows := make([]table.Row, len(data))
	for i := 0; i < len(data); i++ {
		rows[i] = table.SimpleRow{
			data[i].ID,
			data[i].Login,
			data[i].Name,
			data[i].Email,
		}	
	}
	m.table.SetRows(rows)
}

func (m *Model) SetSize(width int, height int) {
	m.table.SetSize(width, height)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.Option = "quit"
		}
	}

	var tablecmd tea.Cmd
	m.table, tablecmd = m.table.Update(msg)
	return m, tea.Batch(cmd, tablecmd)
}

func (m Model) View() string {
	return m.table.View()
}
