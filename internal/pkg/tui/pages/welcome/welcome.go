package welcome

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/messages"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/styles"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/uicontext"
)

const (
	ART = " _               __   _              _                              \n" +
		"| | __   __ _   / _| | | __   __ _  | |  _   _   _ __    ___    ___ \n" +
		"| |/ /  / _` | | |_  | |/ /  / _` | | | | | | | | '_ \\  / __|  / _ \\\n" +
		"|   <  | (_| | |  _| |   <  | (_| | | | | |_| | | |_) | \\__ \\ |  __/\n" +
		"|_|\\_\\  \\__,_| |_|   |_|\\_\\  \\__,_| |_|  \\__, | | .__/  |___/  \\___|\n" +
		"                                         |___/  |_|                 \n" +
		"                                                 \n" +
		"                   Connecting to Kafka cluster...                   \n"
)

type WelcomeStart struct{}
type WelcomeRunning struct{}
type WelcomeDone struct{}

type WelcomeModule struct {
	figString string
	uiCtx     *uicontext.UIContext
	timer     *time.Timer
	active    bool
}

func NewWelcome() WelcomeModule {
	return WelcomeModule{
		figString: ART,
		active:    false,
		uiCtx:     &uicontext.UIContext{},
	}
}

func (w WelcomeModule) Init() tea.Cmd {
	return func() tea.Msg {
		return WelcomeStart{}
	}
}

func (w WelcomeModule) tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		select {
		case <-w.timer.C:
			w.active = false
			return WelcomeDone{}
		default:
			return WelcomeRunning{}
		}
	})
}

func (w WelcomeModule) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case WelcomeStart:
		w.timer = time.NewTimer(750 * time.Millisecond)
		w.active = true
		return w, w.tick()
	case WelcomeRunning:
		if w.active {
			return w, w.tick()
		}
	case WelcomeDone:
		return w, messages.NavigateTo("main")
	case tea.WindowSizeMsg:
		w.uiCtx.Update(msg)
	}
	return w, nil
}

func (w WelcomeModule) View() string {
	var style = styles.WelcomeStyle.
		Width(w.uiCtx.ScreenWidth - 2).
		Height(w.uiCtx.ScreenHeight - 2)

	return style.Render(w.figString)
}
