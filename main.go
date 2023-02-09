package main

import (
	"github.com/snakeice/kafkalypse/internal/pkg/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	if err := a.Run(); err != nil {
		panic(err)
	}

}
