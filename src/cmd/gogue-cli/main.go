package main

import (
	"log"

	gc "github.com/rthornton128/goncurses"
)

const (
	MENU_HEIGHT = 10
	MENU_WIDTH  = 30
)

func main() {
	var active int
	menu := []string{"Choice 1", "Choice 2", "Choice 3", "Choice 4", "Exit"}

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

	my, mx := stdscr.MaxYX()
	y, x := 2, (mx/2)-(MENU_WIDTH/2)

	win, _ := gc.NewWindow(MENU_HEIGHT, MENU_WIDTH, y, x)
	err = win.Keypad(true)
	if err != nil {
		log.Fatal(err)
	}

	stdscr.Print("Use arrow keys to go up and down, Press enter to select")
	stdscr.Refresh()

	printmenu(win, menu, active)

	for {
		ch := stdscr.GetChar()
		switch gc.Key(ch) {
		case 'q':
			return
		case gc.KEY_UP:
			if active == 0 {
				active = len(menu) - 1
			} else {
				active -= 1
			}
		case gc.KEY_DOWN:
			if active == len(menu)-1 {
				active = 0
			} else {
				active += 1
			}
		case gc.KEY_RETURN, gc.KEY_ENTER, gc.Key('\r'):
			stdscr.MovePrintf(my-2, 0, "Choice #%d: %s selected",
				active,
				menu[active])

			err = stdscr.ClearToEOL()
			if err != nil {
				log.Fatal(err)
			}
			stdscr.Refresh()
		default:
			stdscr.MovePrintf(my-2, 0, "Character pressed = %3d/%c",
				ch, ch)

			err = stdscr.ClearToEOL()
			if err != nil {
				log.Fatal(err)
			}
			stdscr.Refresh()
		}

		printmenu(win, menu, active)
	}
}

func printmenu(w *gc.Window, menu []string, active int) {
	y, x := 2, 2
	err := w.Box(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	for i, s := range menu {
		if i == active {
			err := w.AttrOn(gc.A_REVERSE)
			if err != nil {
				log.Fatal(err)
			}

			w.MovePrint(y+i, x, s)

			err = w.AttrOff(gc.A_REVERSE)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			w.MovePrint(y+i, x, s)
		}
	}
	w.Refresh()
}
