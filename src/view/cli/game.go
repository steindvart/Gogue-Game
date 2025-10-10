package cli

import (
	gc "github.com/rthornton128/goncurses"
)

const (
	WINDOW_HEIGHT = 20
	WINDOW_WIDTH  = 50
)

type Game struct {
	W *gc.Window
}

// var (
// 	errMenuAttrOn  = errors.New("cannot set menu attribute on")
// 	errMenuAttrOff = errors.New("cannot set menu attribute off")
// 	errMenuBox     = errors.New("cannot draw menu box")
// )

func NewGame(parent *gc.Window) (*Game, error) {
	// y, x := parent.MaxYX()

	win := parent.Sub(WINDOW_HEIGHT, WINDOW_WIDTH, 0, 0)
	win.Timeout(0)

	// err := win.Keypad(true)
	// if err != nil {
	// 	return nil, err
	// }

	return &Game{W: win}, nil
}

func (m *Game) Render(field [][]int) error {
	m.W.Erase()

	for y, row := range field {
		for x, cell := range row {
			if cell == 1 {
				m.W.MovePrint(y, x, "#")
			} else {
				m.W.MovePrint(y, x, ".")
			}
		}
	}

	m.W.NoutRefresh()
	return nil
}
