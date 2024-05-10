package context

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/pkg/constants"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

type ContextList struct {
	list list.Model
	// keys shortcut.Action

	config *config.Configuration
}

func (m ContextList) Init() tea.Cmd {
	return nil
}

func (m ContextList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.TableStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

	}
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ContextList) View() string {
	h, v := styles.TableStyle.GetFrameSize()
	m.list.SetSize(constants.WindowWidth-h, constants.WindowHeight-v)
	return styles.TableStyle.Render(m.list.View())
}

func NewContextList(config *config.Configuration) *ContextList {
	return &ContextList{
		config: config,
		list:   list.New(ctxToListItem(config), list.NewDefaultDelegate(), 0, 0),
	}
}

type item struct {
	name string
	ctx  *config.Context
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.ctx.BootstrapServers }
func (i item) FilterValue() string { return i.name + i.ctx.BootstrapServers }

func ctxToListItem(config *config.Configuration) []list.Item {
	items := []list.Item{}

	for name, ctx := range config.Contexts {
		items = append(items, item{name: name, ctx: ctx})
	}

	return items
}
