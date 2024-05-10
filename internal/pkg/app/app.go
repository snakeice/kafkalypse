package app

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/pkg/constants"
	"github.com/snakeice/kafkalypse/internal/pkg/kafka"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/connection"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/prompt"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/consumers"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/context"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/topics"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/welcome"
)

type App struct {
	appConfig       *config.Configuration
	kafkaConnection *kafka.Service
	pages           *Pages

	currentPage string
}

func NewApp() *App {
	if len(os.Getenv("DEBUG")) > 0 || true {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	app := App{
		currentPage: "welcome",
		appConfig: &config.Configuration{
			Contexts: map[string]*config.Context{
				"adas": nil,
				"test": nil,
			},
		},
	}

	pages := &Pages{}

	pages.NewPage("Welcome", welcome.NewWelcome("Hello my friend...", nil), false)
	pages.NewPage("Main", nil, false, "main", "home", "m")
	pages.NewPage("Contexts", context.NewContextList(app.appConfig), true, "context", "ctx")
	pages.NewPage("Topics", topics.NewTopics(), true, "topic", "t")
	pages.NewPage("Consumers", consumers.NewConsumers(), true, "consumer", "c")
	pages.NewPage("Brokers", nil, true, "broker", "b")
	pages.NewPage("Producers", nil, true, "producer", "p")
	pages.NewPage("ACLs", nil, true, "acl", "a")

	app.pages = pages
	return &app
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case config.ConfigMsg:
		if msg.Err != nil {
			log.Printf("error: %s", msg.Err)
			return a, tea.Quit
		}

		a.appConfig = msg.Config

		if a.appConfig.CurrentContext == "" {
			cmds = append(cmds, messages.NavigateTo("contexts"))
		} else {
			cmds = append(cmds, messages.NavigateTo("main"),
				kafka.Connect(*a.appConfig.GetCurrentContext()))
		}

		return a, tea.Batch(cmds...)

	case kafka.KafkaConnectionMsg:
		if msg.Err != nil {
			log.Printf("error: %s", msg.Err)
			time.Sleep(5 * time.Second)
			return a, messages.NavigateTo("contexts", true)

		}

		a.kafkaConnection = msg.Svc
		cmds = append(cmds, connection.ConnectionUpdate(connection.ConnectionInfoMsg{
			Brokers:        a.kafkaConnection.BrokersStr(),
			ConectionState: "Connected",
			KafkaVersion:   a.kafkaConnection.Version(),
		}))

	case messages.NavigateToMessage:
		page := a.pages.GetPage(msg.Component)
		a.currentPage = msg.Component
		if page != nil {
			if page.CanAcess || msg.Internal {
				cmds = append(cmds, page.Init())
				return a, tea.Batch(cmds...)
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

	cmd := a.pages.GetPage(a.currentPage).Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(config.LoadConfiguration(), a.pages.GetPage(a.currentPage).Init())
}

func (a *App) View() string {
	return a.pages.GetPage(a.currentPage).View()
}
