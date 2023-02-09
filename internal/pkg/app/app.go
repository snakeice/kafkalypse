package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/pkg/kafka"
	"github.com/snakeice/kafkalypse/internal/pkg/tui"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/prompt"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/container"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/context"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/welcome"
)

type App struct {
	appConfig       *config.Configuration
	kafkaConnection *kafka.KafkaClient
	screen          *tui.UI

	Pages []*Pages
}

func NewApp() (*App, error) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	kafkaConnection, err := kafka.NewKafkaClient()
	if err != nil {
		return nil, err
	}

	screen := tui.NewScreen()

	pages := []*Pages{
		NewPage("Welcome", welcome.NewWelcome(), false),
		NewPage("Main", container.NewContainerModule(), false, "main", "m"),
		NewPage("Contexts", context.NewContextList(appConfig), true, "context", "ctx"),
		NewPage("Topics", screen, true, "topic", "t"),
		NewPage("Consumers", screen, true, "consumer", "c"),
		NewPage("Producers", screen, true, "producer", "p"),
		NewPage("ACLs", screen, true, "acl", "a"),
	}

	return &App{
		appConfig:       appConfig,
		kafkaConnection: kafkaConnection,
		screen:          screen,
		Pages:           pages,
	}, nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case messages.NavigateToMessage:
		for _, page := range a.Pages {
			for _, alias := range page.Aliases {
				if alias == msg.Component {
					return page.Component, nil
				}
			}
		}
	case tea.WindowSizeMsg:
		for _, page := range a.Pages {
			if a.screen.Model != page.Component {
				m, cmd := page.Component.Update(msg)
				cmds = append(cmds, cmd)
				page.Component = m
			}
		}

	case prompt.SubmitMsg:
		if msg.State == prompt.EditingCommand {
			cmds = append(cmds, messages.NavigateTo(msg.Value))
		}
	}

	return nil, tea.Batch(cmds...)
}

func (a *App) Run() error {
	return a.screen.Run(a.Update)
}
