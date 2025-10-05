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
	gc.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	my, mx := stdscr.MaxYX()
	y, x := 2, (mx/2)-(MENU_WIDTH/2)

	win, _ := gc.NewWindow(MENU_HEIGHT, MENU_WIDTH, y, x)
	win.Keypad(true)

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
			stdscr.ClearToEOL()
			stdscr.Refresh()
		default:
			stdscr.MovePrintf(my-2, 0, "Character pressed = %3d/%c",
				ch, ch)
			stdscr.ClearToEOL()
			stdscr.Refresh()
		}

		printmenu(win, menu, active)
	}
}

func printmenu(w *gc.Window, menu []string, active int) {
	y, x := 2, 2
	w.Box(0, 0)
	for i, s := range menu {
		if i == active {
			w.AttrOn(gc.A_REVERSE)
			w.MovePrint(y+i, x, s)
			w.AttrOff(gc.A_REVERSE)
		} else {
			w.MovePrint(y+i, x, s)
		}
	}
	w.Refresh()
}
