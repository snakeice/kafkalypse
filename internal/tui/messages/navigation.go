package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/tools"
)

type NavigateToMessage struct {
	Component string
	Internal  bool
}

func NavigateTo(componentName string, internal ...bool) tea.Cmd {
	isInternal := len(internal) > 0 && internal[0]

	return tools.WrapCmd(NavigateToMessage{
		Component: componentName,
		Internal:  isInternal,
	})
}
