package connection

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/styles"
)

type ConnectionInfoMsg struct {
	Brokers        []string
	ConectionState string
	KafkaVersion   string

	sz messages.SizeMsg
}

type ConnectionModel struct {
	ConnectionInfoMsg
}

func New() ConnectionModel {
	return ConnectionModel{
		ConnectionInfoMsg: ConnectionInfoMsg{
			Brokers:        []string{},
			ConectionState: "Disconnected",
			KafkaVersion:   "Unknown",
		},
	}
}

func (m ConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SizeMsg:
		m.sz = msg

	case ConnectionInfoMsg:
		m.Brokers = msg.Brokers
		m.ConectionState = msg.ConectionState
		m.KafkaVersion = msg.KafkaVersion
	}

	return m, nil
}

func (m ConnectionModel) View() string {
	var components []string

	components = append(components, "Brokers: "+m.getBrokers())
	components = append(components, "Kafka Version: "+m.KafkaVersion)
	components = append(components, m.ConectionState)

	view := lipgloss.JoinVertical(lipgloss.Top, components...)

	return styles.BasicStyle.Render(view)

}

func (m ConnectionModel) Init() tea.Cmd {
	return nil
}

func (m ConnectionModel) getBrokers() string {
	brokers := strings.Join(m.Brokers, ", ")

	wLimit := m.sz.Width - lipgloss.Width("Brokers: ") - 20
	if wLimit < 0 {
		wLimit = 0
	}

	if lipgloss.Width(brokers) < wLimit {
		return brokers
	}

	return brokers[:wLimit] + "..."
}
