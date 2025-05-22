package container

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/tui/components/header"
	"github.com/snakeice/kafkalypse/internal/tui/components/prompt"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/shortcut"
	"github.com/snakeice/kafkalypse/internal/tui/styles"
)

type SetBodyMsg struct {
	Model tea.Model
}

type ContainerModule struct {
	header tea.Model
	prompt tea.Model
	body   tea.Model

	sz messages.SizeMsg
}

func NewContainerModule() ContainerModule {
	return ContainerModule{
		header: header.New(),
		prompt: prompt.New(true),
	}
}

func (m ContainerModule) calculateContentHeight() int {
	height := m.sz.Height
	height -= lipgloss.Height(m.header.View())
	if m.prompt.(prompt.Model).State != prompt.Idle {
		height -= lipgloss.Height(m.prompt.View())
	}
	return height - 4
}

func (m ContainerModule) calculateContentWidth() int {
	return m.sz.Width - 2
}

func (m ContainerModule) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case messages.SizeMsg:
		m.sz = msg
		m.header, _ = m.header.Update(msg)
		m.prompt, _ = m.prompt.Update(msg)
		if m.body != nil {
			m.body, _ = m.body.Update(messages.SizeMsg{
				Width:  m.calculateContentWidth(),
				Height: m.calculateContentHeight(),
			})
		}
		return m, nil
	case SetBodyMsg:
		m.body = msg.Model
		body, cmd := m.body.Update(messages.SizeMsg{
			Width:  m.calculateContentWidth(),
			Height: m.calculateContentHeight(),
		})
		m.body = body
		return m, tea.Batch(cmd)

	case messages.UpdateShortcutsMessage:
		msg.Shortcuts = append(msg.Shortcuts, m.getMainShortcuts()...)
		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)
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

	if view := m.prompt.View(); view != "" {
		components = append(components, view)
	}

	border := "╭"

	sideCount := max((m.sz.Width-4-lipgloss.Width("TOPICS"))/2, 0)

	border += strings.Repeat("─", sideCount)

	border += " " + "TOPICS" + " "

	border += strings.Repeat("─", sideCount)

	border += "╮"

	if m.body != nil {
		v := styles.TableStyle.
			BorderTop(false).
			Width(m.calculateContentWidth()).
			Height(m.calculateContentHeight()).
			Render(m.body.View())

		v = lipgloss.JoinVertical(lipgloss.Top,
			styles.BasicStyle.Render(border),
			v,
		)

		components = append(components, v)
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
