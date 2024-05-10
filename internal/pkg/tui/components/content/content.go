package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/pkg/constants"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ContentModel struct {
	title string
	table tea.Model
}

func New() ContentModel {
	return ContentModel{}
}

func (m ContentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m ContentModel) View() string {

	border := "╭"

	center := (constants.WindowWidth - 2 - lipgloss.Width(m.title)) / 2

	for i := 0; i < center; i++ {
		border += "─"
	}

	border += " " + m.title + " "

	for i := lipgloss.Width(border); i < constants.WindowWidth+1; i++ {
		border += "─"
	}

	border += "╮"

	return lipgloss.JoinVertical(lipgloss.Top,
		styles.BasicStyle.Render(border),

		styles.TableStyle.Copy().
			BorderTop(false).
			Width(constants.WindowWidth).
			Height(constants.WindowHeight).
			Render(m.table.View()),
	)
}

func (m ContentModel) Init() tea.Cmd {
	return m.table.Init()
}

func (m ContentModel) Active() (tea.Model, tea.Cmd) {
	cmd := m.table.Init()
	return m, cmd
}

func (m ContentModel) Deactivate() (tea.Model, tea.Cmd) {
	return m, nil
}
