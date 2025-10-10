package main

import (
	"fmt"

	"gogue/presentation/cli"

	gc "github.com/rthornton128/goncurses"
)

const (
	REQUEIRED_TERMINAL_WIDTH  = 80
	REQUEIRED_TERMINAL_HEIGHT = 50
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
		fmt.Println("Error:", err)
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

	height, width := stdscr.MaxYX()
	if height < REQUEIRED_TERMINAL_HEIGHT || width < REQUEIRED_TERMINAL_WIDTH {
		gc.End()
		return nil, fmt.Errorf("terminal size is too small: need at least %dx%d, got %dx%d",
			REQUEIRED_TERMINAL_WIDTH, REQUEIRED_TERMINAL_HEIGHT, width, height)
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
