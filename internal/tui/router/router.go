package router

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/tui/pages/container"
)

type Router struct {
	routes      map[string]Route
	activeRoute *Route
	container   tea.Model
}

func NewRouter(routes ...Route) *Router {
	r := &Router{
		routes:    make(map[string]Route),
		container: container.NewContainerModule(),
	}
	for _, route := range routes {
		r.routes[route.Key] = route
	}

	return r
}

func (r *Router) Push(model tea.Model) {
	r.activeRoute.stack = append(r.activeRoute.stack, model)
}

func (r *Router) Pop() {
	r.activeRoute.stack = r.activeRoute.stack[:len(r.activeRoute.stack)-1]
}

func (r *Router) GetCurrentStack() []tea.Model {
	if r.activeRoute.stack == nil {
		return nil
	}

	return r.activeRoute.stack
}

func (r *Router) NavigateTo(routeKey string, internal bool) tea.Cmd {
	if route, ok := r.routes[routeKey]; ok {
		if !route.Hidden || internal {
			r.activeRoute = &route

			if route.containered {
				r.container, _ = r.container.Update(container.SetBodyMsg{Model: route.Value})
			} else {
				r.container, _ = r.container.Update(container.SetBodyMsg{Model: nil})
			}

			return r.Init()
		}
	}

	for _, route := range r.routes {
		if slices.Contains(route.Aliases, routeKey) {
			return r.NavigateTo(route.Key, internal)
		}
	}

	return nil
}

func (r *Router) GetCurrentModel() tea.Model {
	if r.activeRoute == nil {
		return nil
	}
	return r.activeRoute.Value
}

func (r *Router) ActivePage() string {
	if r.activeRoute == nil {
		return ""
	}

	return r.activeRoute.Key
}

func (r *Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if r.activeRoute == nil {
		return r, nil
	}

	if r.activeRoute.containered {
		container, cmd := r.container.Update(msg)
		r.container = container
		return r, cmd

	}

	model, cmd := r.activeRoute.Value.Update(msg)
	r.activeRoute.Value = model
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
}

func (r *Router) Init() tea.Cmd {
	if r.activeRoute != nil {
		return r.activeRoute.Value.Init()
	}

	return nil
}

func (r *Router) View() string {
	if r.activeRoute.containered {
		return r.container.View()
	}

	if r.activeRoute != nil {
		return r.activeRoute.Value.View()
	}

	return ""
}
