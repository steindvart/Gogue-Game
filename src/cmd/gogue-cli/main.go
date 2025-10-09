package main

import (
	"log"

	"gogue/presentation/cli"

	gc "github.com/rthornton128/goncurses"
)

func main() {
	stdscr, err := initMainWindowGoncurses()
	if err != nil {
		log.Fatal(err)
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
	if gc.HasColors() {
		err = gc.StartColor()
		if err != nil {
			return nil, err
		}
	}

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
