package router

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Route struct {
	Key         string
	Aliases     []string
	Value       tea.Model
	stack       []tea.Model
	Hidden      bool
	containered bool
}

type option func(*Route)

func NewRoute(key string, value tea.Model, opts ...option) Route {
	r := Route{
		Key:         key,
		Value:       value,
		stack:       []tea.Model{},
		Hidden:      false,
		containered: true,
	}

	for _, opt := range opts {
		opt(&r)
	}

	return r
}

func WithHidden() option {
	return func(r *Route) {
		r.Hidden = true
	}
}

func WithStack(models ...tea.Model) option {
	return func(r *Route) {
		r.stack = models
	}
}

func WithAlias(aliases ...string) option {
	return func(r *Route) {
		r.Aliases = append(r.Aliases, aliases...)
	}
}

func WithContainered() option {
	return func(r *Route) {
		r.containered = true
	}
}

func WithoutContainered() option {
	return func(r *Route) {
		r.containered = false
	}
}

func (r Route) IsContainered() bool {
	return r.containered
}

func (r Route) IsHidden() bool {
	return r.Hidden
}
