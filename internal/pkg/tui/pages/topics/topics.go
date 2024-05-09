package topics

import tea "github.com/charmbracelet/bubbletea"

type Topics struct{}

func NewTopics() *Topics {
	return &Topics{}
}

func (t Topics) Init() tea.Cmd                           { return nil }
func (t Topics) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t Topics) View() string                            { return "TOPICS" }
