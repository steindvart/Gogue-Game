package main

import (
	"fmt"

	"gogue/model/signal"
	"gogue/presentation/cli"
	"gogue/presentation/cli/state"

	"github.com/rivo/tview"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic:", r)
		}
	}()

	app := tview.NewApplication()

	mainMenu := state.NewMainMenu(func(sig signal.Type) {
		switch sig {
		case signal.NewGame:
			// Перейти в состояние игры
			// game.PushState(state.NewGame(...))
		case signal.ShowScoreboard:
			// Показать таблицу лидеров
			// game.PushState(state.NewScoreboard(...))
		case signal.Stop:
			// Завершить приложение
			app.Stop()
		}
	})

	game := &cli.Game{
		States: []state.State{mainMenu},
		App:    app,
	}

	game.Run()
}
