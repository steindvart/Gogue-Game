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
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	gc.Raw(true)
	gc.Echo(false)

	err = gc.Cursor(0)
	if err != nil {
		log.Fatal(err)
	}

	err = stdscr.Clear()
	if err != nil {
		log.Fatal(err)
	}

	err = stdscr.Keypad(true)
	if err != nil {
		log.Fatal(err)
	}

	_, mx := stdscr.MaxYX()
	y, x := 2, (mx/2)-(MENU_WIDTH/2)

	win := stdscr.Sub(MENU_HEIGHT, MENU_WIDTH, y, x)
	stdscr.Timeout(0)
	win.Timeout(0)

	err = win.Keypad(true)
	if err != nil {
		log.Fatal(err)
	}

	game := &cli.Game{
		States: []cli.GameState{cli.NewMainMenuState(win)},
	}

	game.Run()
}
