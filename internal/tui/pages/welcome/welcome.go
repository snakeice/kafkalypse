package welcome

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/constants"
	"github.com/snakeice/kafkalypse/internal/tui/messages"
	"github.com/snakeice/kafkalypse/internal/tui/styles"
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
	msg       string
}

func NewWelcome() WelcomeModule {
	return WelcomeModule{
		figString: ART,
		msg:       "Hello my friend...", //TODO: add a random message here
	}
}

func (w WelcomeModule) Init() tea.Cmd {
	return tea.Tick(1*time.Second, func(time.Time) tea.Msg {
		return WelcomeDone{}
	})
}

func (w WelcomeModule) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(WelcomeDone); ok {
		return w, messages.NavigateTo("main", true)
	}

	return w, nil
}

func (w WelcomeModule) View() string {
	var style = styles.WelcomeStyle.
		Width(constants.WindowWidth - 2).
		Height(constants.WindowHeight - 2)

	return style.Render(fmt.Sprintf(w.figString, w.msg))
}
