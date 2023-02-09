package messages

import tea "github.com/charmbracelet/bubbletea"

// ComponentSizeMessage is a message that is sent when a component is resized.
type ComponentSizeMessage struct {
	Height int
	Width  int
}

type ComponentRefreshMessage struct{}

func ComponentRefresh() tea.Msg {
	return ComponentRefreshMessage{}
}
