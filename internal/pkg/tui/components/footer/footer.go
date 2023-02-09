package footer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FooterModel struct {
}

func New() FooterModel {
	return FooterModel{}
}

func (m FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m FooterModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#999999")).Render(""),

		lipgloss.NewStyle().Foreground(lipgloss.Color("#999999")).Render("Press 'h' to show help"),
	)
}

func (m FooterModel) Init() tea.Cmd {
	return nil
}
