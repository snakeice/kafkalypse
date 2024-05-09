package config

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	NoDefaultContextFound struct{}
	NotFoundAnyContext    struct{}
)

func LoadDefaultContext() tea.Cmd {
	return func() tea.Msg {
		config, err := LoadConfig()
		if err != nil {
			return err
		}

		return config
	}
}
