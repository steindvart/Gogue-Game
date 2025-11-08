package main

import (
	"fmt"
	"gogue/internal/presentation/cli"
	"gogue/internal/presentation/cli/state"

	"github.com/rivo/tview"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic:", r)
		}
	}()

	game := cli.NewGame(tview.NewApplication(), state.NewMainMenu())
	game.Run()
}
