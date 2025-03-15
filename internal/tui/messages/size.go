package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/tools"
)

// SizeMsg is the message to get the size of the terminal
type SizeMsg struct {
	Width  int
	Height int
}

func RefreshSize() tea.Cmd {
	return tools.WrapCmd(RefreshSizeMsg{})
}

type RefreshSizeMsg struct{}
