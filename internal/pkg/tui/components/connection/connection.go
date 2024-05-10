package connection

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/pkg/constants"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ConnectionInfoMsg struct {
	Brokers        []string
	ConectionState string
	KafkaVersion   string
}

func ConnectionUpdate(msg ConnectionInfoMsg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
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
	components = append(components, "Connection State: "+m.ConectionState)
	components = append(components, "Kafka Version: "+m.KafkaVersion)

	view := lipgloss.JoinVertical(lipgloss.Top, components...)

	return styles.BasicStyle.Render(view)

}

func (m ConnectionModel) Init() tea.Cmd {
	return nil
}

func (m ConnectionModel) getBrokers() string {
	brokers := strings.Join(m.Brokers, ", ")

	if lipgloss.Width(brokers) < constants.WindowWidth {
		return brokers
	}

	return brokers[:constants.WindowWidth] + "..."
}
