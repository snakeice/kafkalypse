package app

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/kafka"
	"github.com/snakeice/kafkalypse/internal/tools"
	"github.com/snakeice/kafkalypse/internal/tui/components/connection"
	"github.com/snakeice/kafkalypse/internal/tui/components/prompt"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/pages/consumers"
	"github.com/snakeice/kafkalypse/internal/tui/pages/context"
	"github.com/snakeice/kafkalypse/internal/tui/pages/dbg"
	"github.com/snakeice/kafkalypse/internal/tui/pages/topics"
	"github.com/snakeice/kafkalypse/internal/tui/pages/welcome"
	"github.com/snakeice/kafkalypse/internal/tui/router"
)

type App struct {
	config          *config.Configuration
	kafkaConnection *kafka.Service
	router          *router.Router
	sz              messages.SizeMsg
}

func NewApp() *App {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("error: %s", err)
		return nil
	}

	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	app := App{
		config: cfg,
	}

	// Initialize router with all routes
	rt := router.NewRouter(
		router.NewRoute("main", welcome.NewWelcome(), router.WithHidden()),
		router.NewRoute("context",
			context.NewContextList(app.config),
			router.WithAlias("ctx", "contexts"),
		),
		router.NewRoute("create-context", context.NewCreateContext(app.config)),
		router.NewRoute("topics",
			topics.NewTopics(),
			router.WithAlias("t", "topics"),
		),
		router.NewRoute("consumers", consumers.NewConsumers()),
		router.NewRoute("brokers", nil, router.WithHidden()),
		router.NewRoute("producers", nil, router.WithHidden()),
		router.NewRoute("acls", nil, router.WithHidden()),
		router.NewRoute("dbg", dbg.NewDbg(), router.WithHidden()),
	)

	// Set initial route
	rt.NavigateTo("main", true)
	app.router = rt

	return &app
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case kafka.KafkaConnectionRes:
		if msg.Err != nil {
			log.Printf("error: %s", msg.Err)
			return a, messages.NavigateTo("context", true)
		}

		a.kafkaConnection = msg.KafkaService
		cmds = append(cmds, tools.WrapCmd(connection.ConnectionInfoMsg{
			Brokers:        strings.Split(a.kafkaConnection.BrokersStr(), ","),
			ConectionState: "Connected",
			KafkaVersion:   a.kafkaConnection.Version(),
		}))

		if a.router.ActivePage() == "main" {
			cmds = append(cmds, messages.NavigateTo("topics", true))
		}

	case messages.NavigateToMessage:
		cmds = append(cmds, a.router.NavigateTo(msg.Component, msg.Internal))

	case tea.WindowSizeMsg:
		size := messages.SizeMsg{
			Width:  msg.Width,
			Height: msg.Height,
		}
		a.sz = size
		cmds = append(cmds, tools.WrapCmd(size))

	case messages.RefreshSizeMsg:
		cmds = append(cmds, tools.WrapCmd(a.sz))
	case prompt.SubmitMsg:
		if msg.State == prompt.EditingCommand {
			cmds = append(cmds, messages.NavigateTo(msg.Value))
		}

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	if a.kafkaConnection != nil {
		cmds = append(cmds, a.kafkaConnection.Update(msg))
	}

	_, cmd := a.router.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, a.router.Init())

	if a.config.GetCurrentContext() == nil {
		return messages.NavigateTo("context")
	} else {
		cmds = append(cmds, kafka.GetConn())
		cmds = append(cmds, kafka.Connect(a.config.GetCurrentContext()))
	}
	return tea.Batch(cmds...)
}

func (a *App) View() string {
	return a.router.View()
}
