package main

import (
	"context"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/app"
)

func main() {
	ctx := context.Background()

	f, err := tea.LogToFile("/tmp/kafkalypse.log", "debug")
	if err != nil {
		log.Printf("failed to create log file: %v", err)
		os.Exit(1)
	}
	defer func() { _ = f.Close() }()

	appInstance := app.NewApp()
	if appInstance == nil {
		return
	}

	program := tea.NewProgram(appInstance,
		tea.WithContext(ctx),
		// tea.WithMouseAllMotion(),
		tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("an error occurred: %v\n", err)
	}
}
