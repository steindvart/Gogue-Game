package cli

import (
	"gogue/view/action"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

func Action(k gc.Key) action.Type {
	switch k {
	case gc.KEY_UP:
		return action.MoveUp
	case gc.KEY_DOWN:
		return action.MoveDown
	case gc.KEY_LEFT:
		return action.MoveLeft
	case gc.KEY_RIGHT:
		return action.MoveRight
	case gc.KEY_ENTER:
		fallthrough
	case gc.KEY_RETURN:
		return action.Select
	}

	switch strings.ToLower(string(rune(k))) {
	case "w":
		return action.MoveUp
	case "s":
		return action.MoveDown
	case "a":
		return action.MoveLeft
	case "d":
		return action.MoveRight
	case "e":
		return action.Select
	}
	return action.NoAction
}

func HandleInput(w *gc.Window) action.Type {
	return Action(gc.Key(w.GetChar()))
}
