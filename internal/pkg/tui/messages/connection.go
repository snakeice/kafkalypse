package messages

import tea "github.com/charmbracelet/bubbletea"

// ConnectionMessage is a message that is sent when a connection is established
// or lost.
type ConnectionMessage struct {
	Connected bool
}

func NewConnectionMessage(connected bool) ConnectionMessage {
	return ConnectionMessage{Connected: connected}
}

func ConnectToNewCluster() tea.Msg {
	return ConnectionMessage{Connected: false}
}
