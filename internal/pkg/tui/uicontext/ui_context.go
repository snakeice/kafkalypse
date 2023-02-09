package uicontext

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/kafka"
)

type UIContext struct {
	ScreenHeight int
	ScreenWidth  int

	KafkaConnection *kafka.KafkaClient
}

func (u *UIContext) Update(msg tea.Msg) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		u.ScreenHeight = msg.Height
		u.ScreenWidth = msg.Width
	}
}
