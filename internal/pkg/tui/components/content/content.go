package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/pkg/data"
	"github.com/snakeice/kafkalypse/internal/pkg/kafka"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/table"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ContentModel struct {
	title  string
	uiSize messages.ComponentSizeMessage
	table  tea.Model
}

func New() ContentModel {
	client, err := kafka.NewKafkaClient()
	if err != nil {
		panic(err)
	}

	ds := data.NewKafkaTopics(client)

	return ContentModel{
		title: ds.Name(),
		table: table.New(ds),
	}
}

func (m ContentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(messages.ComponentSizeMessage); ok {
		m.uiSize = msg
	}

	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m ContentModel) View() string {

	border := "╭"

	center := (m.uiSize.Width - 2 - lipgloss.Width(m.title)) / 2

	for i := 0; i < center; i++ {
		border += "─"
	}

	border += " " + m.title + " "

	for i := lipgloss.Width(border); i < m.uiSize.Width+1; i++ {
		border += "─"
	}

	border += "╮"

	return lipgloss.JoinVertical(lipgloss.Top,
		styles.BasicStyle.Render(border),

		styles.TableStyle.Copy().
			BorderTop(false).
			Width(m.uiSize.Width).
			Height(m.uiSize.Height).
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
