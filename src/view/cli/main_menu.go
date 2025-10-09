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
	window *gc.Window
}

var (
	errMenuAttrOn  = errors.New("cannot set menu attribute on")
	errMenuAttrOff = errors.New("cannot set menu attribute off")
	errMenuBox     = errors.New("cannot draw menu box")
)

func NewMainMenu(parent *gc.Window) (*MainMenu, error) {
	_, mx := parent.MaxYX()
	y := 2
	x := (mx / 2) - (MENU_WIDTH / 2)

	win := parent.Sub(MENU_HEIGHT, MENU_WIDTH, y, x)
	win.Timeout(0)

	err := win.Keypad(true)
	if err != nil {
		return nil, err
	}

	return &MainMenu{window: win}, nil
}

func (m *MainMenu) Render(options []string, active int) error {
	err := m.RenderBox()
	if err != nil {
		return err
	}

	err = m.RenderOptions(options, active)
	if err != nil {
		return err
	}

	m.window.Refresh()
	return nil
}

func (m *MainMenu) RenderBox() error {
	err := m.window.Box(0, 0)
	if err != nil {
		return fmt.Errorf("%w: %v", errMenuBox, err)
	}

	return nil
}

func (m *MainMenu) RenderOptions(options []string, active int) error {
	x, y := 2, 2

	for i, s := range options {
		if i == active {
			err := m.window.AttrOn(gc.A_REVERSE)
			if err != nil {
				return fmt.Errorf("%w: %v", errMenuAttrOn, err)
			}

			m.window.MovePrint(y+i, x, s)

			err = m.window.AttrOff(gc.A_REVERSE)
			if err != nil {
				return fmt.Errorf("%w: %v", errMenuAttrOff, err)
			}
		} else {
			m.window.MovePrint(y+i, x, s)
		}
	}

	return nil
}
