package tui

/*
 Draw a screen with modules and template to define the layout.

 The screen is a struct that holds the screen's dimensions and a list of modules.
 The screen is a bubbletea model and can be used as such.

 |--------------------------------------------------|
 |           |                  |                   |
 | <header>  |   <shortcuts>    | <kafkalypse logo> |
 |           |                  |                   |
 |--------------------------------------------------|
 |                                                  |
 |                                                  |
 |                    <module 1>                    |
 |                                                  |
 |                                                  |
 |                                                  |
 |--------------------------------------------------|
 | <Conn info>                                      |
 |--------------------------------------------------|

*/

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/pages/welcome"
)

type UpdateFunc func(tea.Msg) (tea.Model, tea.Cmd)

type UI struct {
	program    *tea.Program
	size       tea.WindowSizeMsg
	Model      tea.Model
	updateFunc UpdateFunc
}

func NewScreen() *UI {
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	ui := &UI{
		Model: welcome.NewWelcome(),
	}

	ui.program = tea.NewProgram(ui, tea.WithAltScreen())
	return ui
}

// Init is called when the screen is initialized.
func (u UI) Init() tea.Cmd {
	return u.Model.Init()
}

// Run starts the screen.
func (u *UI) Run(updateFunc UpdateFunc) error {
	u.updateFunc = updateFunc
	_, err := u.program.Run()
	return err
}

// Update is called every time the screen is updated.
func (u UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return u, tea.Quit
		}

	case tea.WindowSizeMsg:
		u.size = msg
	}

	model, cmd := u.Model.Update(msg)

	u.Model = model

	appModel, appCmd := u.updateFunc(msg)
	if appModel != nil {
		u.Model = appModel
	}

	cmd = tea.Batch(cmd, appCmd)

	return u, cmd
}

// View is called every time the screen is updated.
func (u UI) View() string {
	return u.Model.View()
}
