package main

import (
	"fmt"

	"gogue/presentation/cli"

	gc "github.com/rthornton128/goncurses"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic:", r)
			fmt.Println("probably goncurses error - make sure your terminal supports it and have enough screen size for the game")
		}
	}()

	stdscr, err := initMainWindowGoncurses()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer gc.End()

	game := &cli.Game{
		States: []cli.GameState{cli.NewMainMenu(stdscr)},
		Window: stdscr,
	}

	game.Run()
}

func initMainWindowGoncurses() (*gc.Window, error) {
	stdscr, err := gc.Init()
	if err != nil {
		return nil, err
	}

	gc.Raw(true)
	gc.Echo(false)

	err = gc.Cursor(0)
	if err != nil {
		return nil, err
	}

	err = stdscr.Clear()
	if err != nil {
		return nil, err
	}

	err = stdscr.Keypad(true)
	if err != nil {
		return nil, err
	}
	gc.StdScr().ScrollOk(true)

	stdscr.Timeout(0)

	return stdscr, nil
}
