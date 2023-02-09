package prompt

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type State int

type CancelMsg struct{}
type SubmitMsg struct {
	Value string
	State State
}

const (
	Idle State = iota
	EditingFilter
	EditingCommand
)

type Model struct {
	input    textinput.Model
	State    State
	idleHide bool
	size     tea.WindowSizeMsg
}

func New(idleHide bool) Model {
	ti := textinput.New()
	// ti.BackgroundStyle = styles.TableStyle
	ti.TextStyle.MaxHeight(65)

	m := Model{
		input:    ti,
		idleHide: idleHide,
	}
	m, _ = m.UpdateState(Idle)

	return m
}

func (m Model) UpdateState(state State) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch state {
	case Idle:
		m.input.Prompt = "🦊 > "
		m.input.SetValue("")
		m.input.Blur()
	case EditingFilter:
		m.input.Prompt = "🔍 /"
		cmd := m.input.Focus()
		cmds = append(cmds, cmd)
	case EditingCommand:
		m.input.Prompt = "✨ :"
		cmd := m.input.Focus()
		cmds = append(cmds, cmd)
	}
	m.State = state
	cmds = append(cmds, messages.ComponentRefresh)
	return m, tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			if m.State != Idle {
				break
			}
			var cmd tea.Cmd
			m, cmd = m.UpdateState(EditingCommand)
			cmds = append(cmds, cmd)
		case "/":
			if m.State != Idle {
				break
			}
			var cmd tea.Cmd
			m, cmd = m.UpdateState(EditingFilter)
			cmds = append(cmds, cmd)
		case "esc":
			var cmd tea.Cmd
			m, cmd = m.UpdateState(Idle)
			cmds = append(cmds, cmd)
		case "enter":
			cmds = append(cmds, toCmd(
				SubmitMsg{
					Value: m.Value(),
					State: m.State,
				}))
			var cmd tea.Cmd
			m, cmd = m.UpdateState(Idle)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.size = msg

	case SubmitMsg:
		if msg.Value == "q" && msg.State == EditingCommand {
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.idleHide && m.State == Idle {
		return ""
	}

	view := styles.TableStyle.Copy().
		Width(m.size.Width - 2).
		Render(m.input.View())
	return view
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m Model) SetValue(value string) {
	m.input.SetValue(value)
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func toCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
