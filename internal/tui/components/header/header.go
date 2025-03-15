package header

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/tui/components/connection"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/styles"
)

type Model struct {
	shortcuts      []string
	connectionInfo tea.Model

	sz messages.SizeMsg
}

func New() Model {
	return Model{
		shortcuts:      []string{},
		connectionInfo: connection.New(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.SizeMsg:
		m.sz = msg
		m.connectionInfo, _ = m.connectionInfo.Update(msg)

	case messages.UpdateShortcutsMessage:
		newShortcuts := make([]string, len(msg.Shortcuts))

		for i, shortcut := range msg.Shortcuts {
			newShortcuts[i] = fmt.Sprintf("%s %s", shortcut.Shortcuts[0], shortcut.Description)
		}

		m.shortcuts = newShortcuts

		return m, nil
	}

	connInfo, cmd := m.connectionInfo.Update(msg)
	m.connectionInfo = connInfo

	return m, cmd
}

func (m Model) Init() tea.Cmd {
	return nil // nothing to do here
}

func (m Model) View() string {
	headerStyle := styles.HeaderComponent.Height(3).Width(m.sz.Width)

	shortcuts := strings.Join(m.shortcuts, "\n")
	connectionInfo := m.connectionInfo.View()

	// Calculate remaining width for shortcuts
	shortcutsWidth := m.sz.Width - lipgloss.Width(connectionInfo)

	view := lipgloss.JoinHorizontal(lipgloss.Top,
		connectionInfo,
		headerStyle.Width(shortcutsWidth).Render(shortcuts),
	)

	return styles.BasicStyle.Width(m.sz.Width).Render(view)
}
