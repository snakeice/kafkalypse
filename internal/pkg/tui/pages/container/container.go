package container

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/pkg/constants"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/header"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/prompt"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/shortcut"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ContainerModule struct {
	header tea.Model
	prompt tea.Model
}

func NewContainerModule() ContainerModule {
	return ContainerModule{
		header: header.New(),

		prompt: prompt.New(true),
	}
}

func (m ContainerModule) calculateContentHeight() int {
	height := constants.WindowHeight
	height -= lipgloss.Height(m.header.View())
	if m.prompt.(prompt.Model).State != prompt.Idle {
		height -= lipgloss.Height(m.prompt.View())
	}
	// height -= lipgloss.Height(m.footer.View())
	return height
}

func (m ContainerModule) tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return messages.ComponentRefreshMessage{}
	})
}

func (m ContainerModule) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case messages.ComponentRefreshMessage:

		// cmds = append(cmds, cmd, m.tick())

		return m, tea.Batch(cmds...)

	case messages.UpdateShortcutsMessage:
		msg.Shortcuts = append(msg.Shortcuts, m.getMainShortcuts()...)

		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	}

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.prompt, cmd = m.prompt.Update(msg)
	cmds = append(cmds, cmd)

	// m.footer, cmd = m.footer.Update(msg)
	// cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ContainerModule) View() string {
	var components []string

	components = append(components, m.header.View())

	if view := m.prompt.View(); view != "" || true {
		components = append(components, view)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		components...,
	)

	return styles.BasicStyle.Render(view)
}

func (m ContainerModule) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		messages.ComponentRefresh,
		messages.UpdateShortcuts,
	)
}

func (m ContainerModule) getMainShortcuts() []shortcut.Action {
	return []shortcut.Action{
		{
			Description: "Quit",
			Shortcuts: []string{
				":q",
				"ctrl+c",
			},
			Action: "Quit",
		},
		{
			Description: "Help",
			Shortcuts: []string{
				":h",
			},
			Action: "Help",
		},
	}
}

func (m ContainerModule) Active() (tea.Model, tea.Cmd) {
	return m, m.Init()
}

func (m ContainerModule) Deactivate() (tea.Model, tea.Cmd) {
	return m, nil
}
