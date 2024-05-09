package welcome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/constants"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
)

const (
	ART = " _               __   _              _                              \n" +
		"| | __   __ _   / _| | | __   __ _  | |  _   _   _ __    ___    ___ \n" +
		"| |/ /  / _` | | |_  | |/ /  / _` | | | | | | | | '_ \\  / __|  / _ \\\n" +
		"|   <  | (_| | |  _| |   <  | (_| | | | | |_| | | |_) | \\__ \\ |  __/\n" +
		"|_|\\_\\  \\__,_| |_|   |_|\\_\\  \\__,_| |_|  \\__, | | .__/  |___/  \\___|\n" +
		"                                         |___/  |_|                 \n" +
		"                                                 \n" +
		"%s\n"
)

type WelcomeDone struct{}

type WelcomeModule struct {
	figString string
}

func NewWelcome(msg string, returnMsg tea.Msg) WelcomeModule {
	return WelcomeModule{
		figString: ART,
	}
}

func (w WelcomeModule) Init() tea.Cmd {
	return nil
}

func (w WelcomeModule) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case WelcomeDone:
		return w, messages.NavigateTo("main", true)
	}
	return w, nil
}

func (w WelcomeModule) View() string {
	var style = styles.WelcomeStyle.
		Width(constants.WindowWidth - 2).
		Height(constants.WindowHeight - 2)

	return style.Render(w.figString)
}
