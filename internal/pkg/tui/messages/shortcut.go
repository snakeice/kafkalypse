package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/shortcut"
)

// UpdateShortcutsMessage is a message that is sent when the shortcuts are
// updated.
type UpdateShortcutsMessage struct {
	Shortcuts []shortcut.Action
}

func UpdateShortcuts() tea.Msg {
	return UpdateShortcutsMessage{
		Shortcuts: []shortcut.Action{},
	}
}
