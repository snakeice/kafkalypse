package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/pkg/app"
)

func main() {
	program := tea.NewProgram(app.NewApp(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}

}
