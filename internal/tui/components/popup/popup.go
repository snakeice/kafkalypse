package popup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/styles"
)

type SubmitMsg struct {
	Values map[string]string
}

type Field struct {
	Name        string
	Value       string
	Placeholder string
}

type Model struct {
	title       string
	fields      []Field
	inputs      []textinput.Model
	activeInput int
	visible     bool
	width       int
	height      int

	sz messages.SizeMsg
}

func New(title string, fields []Field) Model {
	inputs := make([]textinput.Model, len(fields))
	for i := range fields {
		input := textinput.New()
		input.Placeholder = fields[i].Placeholder
		input.Width = 30
		input.SetValue(fields[i].Value)
		inputs[i] = input
	}

	if len(inputs) > 0 {
		inputs[0].Focus()
	}

	return Model{
		title:       title,
		fields:      fields,
		inputs:      inputs,
		activeInput: 0,
		visible:     false,
		width:       40,
		height:      len(fields)*3 + 6,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case messages.SizeMsg:
		m.sz = msg
		for i := range m.inputs {
			m.inputs[i].Width = m.sz.Width - 10
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.visible = false
			return m, nil
		case "enter":
			if m.activeInput == len(m.inputs)-1 {
				// Submit form
				values := make(map[string]string)
				for i, field := range m.fields {
					values[field.Name] = m.inputs[i].Value()
				}
				m.visible = false
				return m, func() tea.Msg {
					return SubmitMsg{Values: values}
				}
			}
			// Move to next field
			m.activeInput = (m.activeInput + 1) % len(m.inputs)
			for i := range m.inputs {
				if i == m.activeInput {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
		case "tab":
			m.activeInput = (m.activeInput + 1) % len(m.inputs)
			for i := range m.inputs {
				if i == m.activeInput {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
		case "shift+tab":
			m.activeInput--
			if m.activeInput < 0 {
				m.activeInput = len(m.inputs) - 1
			}
			for i := range m.inputs {
				if i == m.activeInput {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
		}
	}

	// Handle character input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for i := range m.inputs {
		if i == m.activeInput {
			newInput, cmd := m.inputs[i].Update(msg)
			m.inputs[i] = newInput
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Calculate popup position
	screenWidth := m.sz.Width
	screenHeight := m.sz.Height

	// Create popup content
	var content string
	content += lipgloss.NewStyle().Bold(true).Render(m.title) + "\n\n"

	// Add form fields
	for i, field := range m.fields {
		content += lipgloss.NewStyle().Bold(true).Render(field.Name+":") + "\n"
		content += m.inputs[i].View() + "\n"
	}

	content += "\n" + lipgloss.NewStyle().
		Faint(true).
		Render("Press Enter to submit • Esc to cancel • Tab to navigate")

	// Style the popup
	popup := styles.BasicStyle.
		Width(m.width).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Render(content)

	// Position the popup
	return lipgloss.Place(
		screenWidth,
		screenHeight,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)
}

func (m Model) IsVisible() bool {
	return m.visible
}

func (m *Model) Hide() {
	m.visible = false
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func (m *Model) Show() {
	m.visible = true
	m.activeInput = 0
	for i := range m.inputs {
		if i == 0 {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	// Reset input values
	for i, field := range m.fields {
		m.inputs[i].SetValue(field.Value)
	}
}
