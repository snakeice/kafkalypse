package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type PageModel interface {
	tea.Model
}

type Pages []*Page
type Page struct {
	Title     string
	CanAcess  bool
	Aliases   []string
	Component tea.Model
}

func NewPage(title string, model PageModel, canAcess bool, aliases ...string) *Page {
	return &Page{
		Title:     title,
		Aliases:   aliases,
		Component: model,
		CanAcess:  canAcess,
	}
}

func (p *Pages) GetPage(n string) *Page {
	n = strings.ToLower(n)
	for _, page := range *p {
		if strings.ToLower(page.Title) == n {
			return page
		}
		for _, alias := range page.Aliases {

			if alias == n {
				return page
			}
		}
	}
	return nil
}

func (p *Page) Update(msg tea.Msg) tea.Cmd {
	m, cmd := p.Component.Update(msg)

	p.Component = m

	return cmd
}
