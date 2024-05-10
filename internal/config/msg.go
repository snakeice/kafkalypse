package config

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ConfigMsg struct {
	Config *Configuration
	Err    error
}

func LoadConfiguration() tea.Cmd {
	return func() tea.Msg {
		config, err := LoadConfig()
		if err != nil {
			return ConfigMsg{nil, err}
		}

		return ConfigMsg{config, nil}
	}
}
