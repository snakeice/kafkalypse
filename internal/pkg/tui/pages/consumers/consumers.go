package consumers

import tea "github.com/charmbracelet/bubbletea"

type Consumers struct {
}

func NewConsumers() *Consumers {
	return &Consumers{}
}

func (c Consumers) Init() tea.Cmd                           { return nil }
func (c Consumers) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return c, nil }
func (c Consumers) View() string                            { return "Consumers" }
