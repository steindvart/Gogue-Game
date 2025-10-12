package main

import (
	"fmt"

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

	game := cli.NewGame(tview.NewApplication(), state.NewMainMenu())
	game.Run()
}
