package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/container"
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
		Component: container.NewContainerModule(model),
		CanAcess:  canAcess,
	}
}

func (p *Pages) NewPage(title string, model PageModel, canAcess bool, aliases ...string) {
	*p = append(*p, NewPage(title, model, canAcess, aliases...))
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

func (p *Page) Init() tea.Cmd {
	return p.Component.Init()
}

func (p *Page) View() string {
	return p.Component.View()
}
