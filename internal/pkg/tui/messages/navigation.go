package messages

import tea "github.com/charmbracelet/bubbletea"

type NavigateToMessage struct {
	Component string
	Internal  bool
}

func NavigateTo(componentName string, internal ...bool) tea.Cmd {
	return func() tea.Msg {
		isInternal := false
		if len(internal) > 0 {
			isInternal = internal[0]
		}

		return NavigateToMessage{
			Component: componentName,
			Internal:  isInternal,
		}
	}
}
