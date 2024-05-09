package connection

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ConnectionModel struct {
	Brokers        []string
	ConectionState string
	KafkaVersion   string

	size tea.WindowSizeMsg
}

func New() ConnectionModel {
	return ConnectionModel{
		Brokers:        []string{},
		ConectionState: "Disconnected",
		KafkaVersion:   "Unknown",
	}
}

func (m ConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	if lipgloss.Width(brokers) < m.size.Width {
		return brokers
	}

	return brokers[:m.size.Width] + "..."
}
