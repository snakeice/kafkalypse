package pagesstack

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Page interface {
	tea.Model
	Active() (tea.Model, tea.Cmd)
	Deactivate() (tea.Model, tea.Cmd)
}
