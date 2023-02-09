package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/shortcut"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ColHead struct {
	Name string
	Perc float64
}

type datasource interface {
	Len() int
	At(int) []string
	Shortcuts() []shortcut.Action
	Cols() []ColHead
}

type KeyMap struct {
	LineUp       key.Binding
	LineDown     key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	GotoTop      key.Binding
	GotoBottom   key.Binding
}

type Model struct {
	datasource datasource

	viewport viewport.Model
	focus    bool
	cursor   int

	uiSize messages.ComponentSizeMessage

	KeyMap KeyMap

	start int
	end   int
}

func DefaultKeyMap() KeyMap {
	const spacebar = " "
	return KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp("b/pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", spacebar),
			key.WithHelp("f/pgdn", "page down"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp("u", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp("d", "½ page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
	}
}

func New(datasource datasource) Model {
	return Model{
		viewport:   viewport.New(0, 10),
		datasource: datasource,
		cursor:     0,
		KeyMap:     DefaultKeyMap(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case messages.ComponentSizeMessage:
		if m.uiSize == msg {
			return m, nil
		}

		m.uiSize = msg

		m.viewport.Height = m.uiSize.Height - 1 // 3 for the header
		m.viewport.Width = m.uiSize.Width

		m.updateViewport()
	case tea.KeyMsg:
		// if !m.focus {
		// 	return m, nil
		// }

		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.PageUp):
			m.MoveUp(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) MoveUp(n int) {
	m.cursor = clamp(m.cursor-n, 0, m.datasource.Len()-1)

	switch {
	case m.start == 0:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset, 0, m.cursor))
	case m.cursor < m.viewport.Height:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset+n, 0, m.cursor))
	case m.viewport.YOffset >= 1:
		m.viewport.YOffset = clamp(m.viewport.YOffset+n, 1, m.cursor)
	}

	m.updateViewport()
}

func (m *Model) MoveDown(n int) {
	m.cursor = clamp(m.cursor+n, 0, m.datasource.Len()-1)
	m.updateViewport()

	switch {
	case m.end == m.datasource.Len():
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.viewport.Height))
	case m.cursor > (m.end-m.start)/2:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.cursor))
	case m.viewport.YOffset > 1:
	case m.cursor > m.viewport.YOffset+m.viewport.Height-1:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset+1, 0, 1))
	}
}

func (m *Model) GotoTop() {
	m.MoveUp(m.cursor)
}

func (m *Model) GotoBottom() {
	m.MoveDown(m.datasource.Len())
}

func (m *Model) headersView() string {
	var s = make([]string, len(m.datasource.Cols()))

	for i, col := range m.datasource.Cols() {
		width := int(float64(m.uiSize.Width) * col.Perc)

		style := lipgloss.NewStyle().
			Width(width).
			MaxWidth(width).
			Inline(true)

		s[i] = style.Render(runewidth.Truncate(col.Name, width, "..."))
	}

	return styles.TableHeader.
		Width(m.uiSize.Width).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, s...))
}

func (m *Model) renderRow(index int) string {
	var s = make([]string, len(m.datasource.Cols()))

	rowData := m.datasource.At(index)

	for i, col := range m.datasource.Cols() {
		width := int(float64(m.uiSize.Width) * col.Perc)

		style := lipgloss.NewStyle().
			Width(width).
			MaxWidth(width).
			Inline(true)

		if index == m.cursor {
			style = styles.TableHeaderSelected.Copy().Width(width).MaxWidth(width)
		}

		s[i] = style.Render(runewidth.Truncate(rowData[i], width, "..."))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *Model) updateViewport() {
	rows := make([]string, 0, m.datasource.Len())

	if m.cursor >= 0 {
		m.start = clamp(m.cursor-m.viewport.Height, 0, m.cursor)
	} else {
		m.start = 0
	}

	m.end = clamp(m.cursor+m.viewport.Height, m.cursor, m.datasource.Len())

	for i := m.start; i < m.end; i++ {
		rows = append(rows, m.renderRow(i))
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)
}

func (m Model) View() string {
	m.updateViewport()
	return m.headersView() + "\n" + m.viewport.View()
}

func (m Model) Init() tea.Cmd { return nil }

func (m *Model) Focus() {
	m.focus = true
	m.updateViewport()
}

func (m *Model) Blur() {
	m.focus = false
	m.updateViewport()
}
