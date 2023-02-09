package messages

import tea "github.com/charmbracelet/bubbletea"

type NavigateToMessage struct {
	Component string
}

func NavigateTo(componentName string) tea.Cmd {
	return func() tea.Msg {
		return NavigateToMessage{
			Component: componentName,
		}
	}
}
