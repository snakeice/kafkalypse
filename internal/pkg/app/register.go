package app

import tea "github.com/charmbracelet/bubbletea"

type Pages struct {
	Title     string
	CanAcess  bool
	Aliases   []string
	Component tea.Model
}

func NewPage(title string, model tea.Model, canAcess bool, aliases ...string) *Pages {
	return &Pages{
		Title:     title,
		Aliases:   aliases,
		Component: model,
		CanAcess:  canAcess,
	}
}
