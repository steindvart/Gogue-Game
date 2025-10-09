package main

import (
	"log"

	"gogue/presentation/cli"

	gc "github.com/rthornton128/goncurses"
)

const (
	MENU_HEIGHT = 10
	MENU_WIDTH  = 30
)

func main() {
	stdscr, err := initMainWindowGoncurses()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	_, mx := stdscr.MaxYX()
	y, x := 2, (mx/2)-(MENU_WIDTH/2)

	win := stdscr.Sub(MENU_HEIGHT, MENU_WIDTH, y, x)
	win.Timeout(0)

	err = win.Keypad(true)
	if err != nil {
		log.Fatal(err)
	}

	game := &cli.Game{
		States: []cli.GameState{cli.NewMainMenu(win)},
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
