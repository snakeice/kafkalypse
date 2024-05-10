package container

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/header"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/prompt"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/shortcut"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ContainerModule struct {
	header tea.Model
	prompt tea.Model
	body   tea.Model
}

func NewContainerModule(body tea.Model) ContainerModule {
	return ContainerModule{
		header: header.New(),
		prompt: prompt.New(true),
		body:   body,
	}
}

// func (m ContainerModule) calculateContentHeight() int {
// 	height := constants.WindowHeight
// 	height -= lipgloss.Height(m.header.View())
// 	if m.prompt.(prompt.Model).State != prompt.Idle {
// 		height -= lipgloss.Height(m.prompt.View())
// 	}
// 	// height -= lipgloss.Height(m.footer.View())
// 	return height
// }

func (m ContainerModule) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd

	switch msg := msg.(type) {
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

	if m.body != nil {
		m.body, cmd = m.body.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ContainerModule) View() string {
	var components []string

	components = append(components, m.header.View())

	if view := m.prompt.View(); view != "" || true {
		components = append(components, view)
	}

	if m.body != nil {
		components = append(components, m.body.View())
	}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		components...,
	)

	return styles.BasicStyle.Render(view)
}

func (m ContainerModule) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.header.Init(),
		m.prompt.Init(),
		textinput.Blink,
		messages.UpdateShortcuts,
	}

	if m.body != nil {
		cmds = append(cmds, m.body.Init())
	}

	return tea.Batch(cmds...)
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
