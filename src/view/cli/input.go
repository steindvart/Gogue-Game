package cli

import (
	"gogue/view"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

func Action(k gc.Key) view.ActionType {
	switch k {
	case gc.KEY_UP:
		return view.MoveUp
	case gc.KEY_DOWN:
		return view.MoveDown
	case gc.KEY_LEFT:
		return view.MoveLeft
	case gc.KEY_RIGHT:
		return view.MoveRight
	case gc.KEY_ENTER:
		return view.Select
	}

	switch strings.ToLower(string(rune(k))) {
	case "w":
		return view.MoveUp
	case "s":
		return view.MoveDown
	case "a":
		return view.MoveLeft
	case "d":
		return view.MoveRight
	case "e":
		return view.Select
	}
	return view.NoAction
}
