package cli

import (
	"errors"
	"fmt"

	gc "github.com/rthornton128/goncurses"
)

const (
	MENU_HEIGHT = 10
	MENU_WIDTH  = 30
)

type MainMenu struct {
	W *gc.Window
}

var (
	errMenuAttrOn  = errors.New("cannot set menu attribute on")
	errMenuAttrOff = errors.New("cannot set menu attribute off")
	errMenuBox     = errors.New("cannot draw menu box")
)

func NewMainMenu(parent *gc.Window) (*MainMenu, error) {
	_, mx := parent.MaxYX()
	y, x := 2, (mx/2)-(MENU_WIDTH/2)

	win := parent.Sub(MENU_HEIGHT, MENU_WIDTH, y, x)
	win.Timeout(0)

	err := win.Keypad(true)
	if err != nil {
		return nil, err
	}

	return &MainMenu{W: win}, nil
}

func (m *MainMenu) Render(options []string, active int) error {
	m.W.Erase()

	x, y := 2, 2
	err := m.W.Box(0, 0)
	if err != nil {
		return fmt.Errorf("%w: %v", errMenuBox, err)
	}

	for i, s := range options {
		if i == active {
			err := m.W.AttrOn(gc.A_REVERSE)
			if err != nil {
				return fmt.Errorf("%w: %v", errMenuAttrOn, err)
			}

			m.W.MovePrint(y+i, x, s)

			err = m.W.AttrOff(gc.A_REVERSE)
			if err != nil {
				return fmt.Errorf("%w: %v", errMenuAttrOff, err)
			}
		} else {
			m.W.MovePrint(y+i, x, s)
		}
	}

	m.W.NoutRefresh()
	return nil
}
