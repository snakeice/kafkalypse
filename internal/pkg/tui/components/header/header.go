package header

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/connection"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type Model struct {
	shortcuts      []string
	connectionInfo connection.ConnectionModel
}

func New() Model {
	return Model{
		shortcuts:      []string{},
		connectionInfo: connection.New(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.UpdateShortcutsMessage:
		newShortcuts := make([]string, len(msg.Shortcuts))

		for i, shortcut := range msg.Shortcuts {
			newShortcuts[i] = fmt.Sprintf("%s %s", shortcut.Shortcuts[0], shortcut.Description)
		}

		m.shortcuts = newShortcuts

		return m, nil
	}

	connInfo, cmd := m.connectionInfo.Update(msg)
	m.connectionInfo = connInfo.(connection.ConnectionModel)

	return m, cmd

}

func (m Model) Init() tea.Cmd {
	return nil // nothing to do here
}

func (m Model) View() string {
	styles.HeaderComponent.Height(3)

	shortcuts := strings.Join(m.shortcuts, "\n")

	view := lipgloss.JoinHorizontal(lipgloss.Right,
		m.connectionInfo.View(),
		styles.HeaderComponent.Render(shortcuts),
	)

	return view
}
