package kafka

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/config"
)

func Connect(params config.Context) tea.Cmd {
	return func() tea.Msg {

		msg := KafkaConnectionMsg{}

		msg.Svc, msg.Err = NewKafkaClient()

		return msg
	}
}

// TODO
// - Topic list
// - Topic detail
// - Consumer group list
// - Consumer group detail
// - Consumer group offset reset
// - Consume topic messages
// - Produce topic messages
// - Topic creation
// - Topic deletion
// - Topic configuration
// - Schema registry?
