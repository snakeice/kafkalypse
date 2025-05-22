package shortcuts

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Action struct {
	Key         string
	Description string
	Handler     func() tea.Cmd
}

type Shortcuts struct {
	actions map[string]Action
	focus   bool
}

func New() *Shortcuts {
	return &Shortcuts{
		actions: make(map[string]Action),
		focus:   true,
	}
}

func (s *Shortcuts) Add(action Action) {
	s.actions[action.Key] = action
}

func (s *Shortcuts) Handle(msg tea.KeyMsg) tea.Cmd {
	if !s.focus {
		return nil
	}
	if action, ok := s.actions[msg.String()]; ok {
		return action.Handler()
	}
	return nil
}

func (s *Shortcuts) GetHelpText() string {
	var descriptions []string
	for _, action := range s.actions {
		descriptions = append(descriptions, action.Description)
	}
	return "• " + joinStrings(descriptions, " • ")
}

func (s *Shortcuts) Focus() {
	s.focus = true
}

func (s *Shortcuts) Blur() {
	s.focus = false
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
