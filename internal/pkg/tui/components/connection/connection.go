package connection

import tea "github.com/charmbracelet/bubbletea"

type ConnectionModel struct {
}

func New() ConnectionModel {
	return ConnectionModel{}
}

func (m ConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConnectionModel) View() string {
	return "X"
}

func (m ConnectionModel) Init() tea.Cmd {
	return nil
}
