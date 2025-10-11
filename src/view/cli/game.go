package cli

import (
	"errors"
	"fmt"

	gc "github.com/rthornton128/goncurses"
)

type Game struct {
	window *gc.Window
}

var (
	errBox = errors.New("cannot draw game window box")
)

func NewGame(parent *gc.Window, fieldWidth, fieldHeight int) (*Game, error) {
	_, mx := parent.MaxYX()

	y := 20
	x := (mx / 2) - (fieldWidth / 2)

	win := parent.Sub(fieldHeight, fieldWidth, y, x)

	return &Game{window: win}, nil
}

func (g *Game) Render(field [][]int) error {
	g.window.Erase()

	err := g.window.Box(0, 0)
	if err != nil {
		return fmt.Errorf("%w: %v", errBox, err)
	}

	g.drawField(field)

	g.window.NoutRefresh()
	return nil
}

func (g *Game) drawField(field [][]int) {
	offset := 1 // смещение из-за границ

	for y, row := range field {
		for x, cell := range row {
			if cell == 1 {
				g.window.MovePrint(y+offset, x+offset, "@")
			} else {
				g.window.MovePrint(y+offset, x+offset, " ")
			}
		}
	}
}
