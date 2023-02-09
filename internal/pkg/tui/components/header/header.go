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

const (
	Logo = `╭╮╭━╮╱╱╱╭━┳╮╱╱╱╱╭╮
┃┃┃╭╯╱╱╱┃╭┫┃╱╱╱╱┃┃
┃╰╯╯╭━━┳╯╰┫┃╭┳━━┫┃╭╮╱╭┳━━┳━━┳━━╮
┃╭╮┃┃╭╮┣╮╭┫╰╯┫╭╮┃┃┃┃╱┃┃╭╮┃━━┫┃━┫
┃┃┃╰┫╭╮┃┃┃┃╭╮┫╭╮┃╰┫╰━╯┃╰╯┣━━┃┃━┫
╰╯╰━┻╯╰╯╰╯╰╯╰┻╯╰┻━┻━╮╭┫╭━┻━━┻━━╯
╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╰━━╯╰╯ %s`
)

type Model struct {
	logo           string
	version        string
	shortcuts      []string
	connectionInfo connection.ConnectionModel
	size           tea.WindowSizeMsg
}

func New() Model {
	return Model{
		logo:           Logo,
		version:        "v0.0.1",
		shortcuts:      []string{},
		connectionInfo: connection.New(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg

	case messages.UpdateShortcutsMessage:
		newShortcuts := make([]string, len(msg.Shortcuts))

		for i, shortcut := range msg.Shortcuts {
			newShortcuts[i] = fmt.Sprintf("%s %s", shortcut.Shortcuts[0], shortcut.Description)
		}

		m.shortcuts = newShortcuts

		return m, messages.ComponentRefresh
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil // nothing to do here
}

func (m Model) View() string {
	styles.HeaderComponent.Height(lipgloss.Height(Logo) + 1)

	shortcuts := strings.Join(m.shortcuts, "\n")
	logo := fmt.Sprintf(m.logo, m.version)

	view := lipgloss.JoinHorizontal(lipgloss.Right,
		styles.HeaderComponent.Render(logo),
		styles.HeaderComponent.Render(shortcuts),
	)

	style := styles.HeaderComponent.Copy().Width(m.size.Width - lipgloss.Width(view)).AlignHorizontal(lipgloss.Right)

	view = lipgloss.JoinHorizontal(lipgloss.Right,
		view,
		style.Render(m.connectionInfo.View()),
	)

	return view
}
