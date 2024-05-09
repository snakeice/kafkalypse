package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/pkg/constants"
	"github.com/snakeice/kafkalypse/internal/pkg/kafka"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/prompt"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/container"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/context"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/welcome"
)

type App struct {
	appConfig       *config.Configuration
	kafkaConnection *kafka.Service
	pages           Pages
}

func NewApp() *App {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	pages := Pages{
		NewPage("Welcome", welcome.NewWelcome("Hello", nil), false),
		NewPage("Main", container.NewContainerModule(), false, "main", "home", "m"),
		NewPage("Contexts", context.NewContextList(&config.Configuration{}), true, "context", "ctx"),
		// NewPage("Topics", topics.NewTopics(), true, "topic", "t"),
		// NewPage("Consumers", consumers.NewConsumers(), true, "consumer", "c"),
		NewPage("Brokers", nil, true, "broker", "b"),
		NewPage("Producers", nil, true, "producer", "p"),
		NewPage("ACLs", nil, true, "acl", "a"),
	}

	return &App{
		pages: pages,
	}
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case config.Configuration:
		a.appConfig = &msg
		return a, nil
	case messages.NavigateToMessage:
		page := a.pages.GetPage(msg.Component)
		if page != nil {
			if page.CanAcess || msg.Internal {
				return nil, tea.Batch(cmds...)
			}
		}

	case tea.WindowSizeMsg:
		constants.WindowHeight = msg.Height
		constants.WindowWidth = msg.Width

	case prompt.SubmitMsg:
		if msg.State == prompt.EditingCommand {
			cmds = append(cmds, messages.NavigateTo(msg.Value))
		}

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	cmd := a.pages.GetPage("main").Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(config.LoadDefaultContext(), a.pages.GetPage("main").Component.Init())
}

func (a *App) View() string {
	return a.pages.GetPage("main").Component.View()
}
