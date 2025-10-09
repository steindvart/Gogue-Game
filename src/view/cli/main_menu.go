package cli

import (
	"errors"
	"fmt"
	"time"

	gc "github.com/rthornton128/goncurses"
)

const (
	MENU_HEIGHT = 30
	MENU_WIDTH  = 70
)

// ASCII-art title
var titleArt = []string{
	"   _____                           _____                      ",
	"  / ____|                         / ____|                     ",
	" | |  __  ___   __ _ _   _  ___  | |  __  __ _ _ __ ___   ___ ",
	" | | |_ |/ _ \\ / _` | | | |/ _ \\ | | |_ |/ _` | '_ ` _ \\ / _ \\",
	" | |__| | (_) | (_| | |_| |  __/ | |__| | (_| | | | | | |  __/",
	"  \\_____|\\___/ \\__, |\\__,_|\\___|  \\_____|\\__,_|_| |_| |_|\\___|",
	"                __/ |                                         ",
	"               |___/                                          ",
}

type MainMenu struct {
	window *gc.Window
	colors []int
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

	var colors []int
	if gc.HasColors() {
		colors = createRainbowPairs()
	}

	return &MainMenu{window: win, colors: colors}, nil
}

func createRainbowPairs() []int {
	var colors = []int{
		gc.C_RED, gc.C_YELLOW, gc.C_GREEN, gc.C_CYAN, gc.C_BLUE, gc.C_MAGENTA,
	}

	for i, color := range colors {
		gc.InitPair(int16(i+1), int16(color), gc.C_BLACK)
	}

	return colors
}

func (m *MainMenu) Render(options []string, active int) error {
	// err := m.RenderBox()
	// if err != nil {
	// 	return err
	// }

	m.RenderTitle()

	err := m.RenderOptions(options, active)
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
	_, maxX := m.window.MaxYX()
	y := 10
	x := (maxX / 2) - 6

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

func (m *MainMenu) RenderTitle() {
	if gc.HasColors() {
		// Для плавности: используем текущее время как frame
		frame := int(time.Now().UnixNano() / 75000000) // ~15 кадров в сек
		drawGradientTitle(m.window, frame, m.colors)
	} else {
		drawRawTitle(m.window)
	}
}

func drawGradientTitle(win *gc.Window, frame int, colors []int) {
	_, maxX := win.MaxYX()
	titleWidth := len(titleArt[0])
	startY := 1                       // чуть ниже начала окна
	startX := (maxX - titleWidth) / 2 // по центру с учётом ширины

	for row, line := range titleArt {
		for col, ch := range line {
			colorIdx := (col + frame) % len(colors)
			win.AttrOn(gc.ColorPair(int16(colorIdx + 1)))
			win.MovePrint(startY+row, startX+col, string(ch))
			win.AttrOff(gc.ColorPair(int16(colorIdx + 1)))
		}
	}
}

func drawRawTitle(win *gc.Window) {
	_, maxX := win.MaxYX()
	titleWidth := len(titleArt[0])
	startY := 1                       // чуть ниже начала окна
	startX := (maxX - titleWidth) / 2 // по центру с учётом ширины

	for row, line := range titleArt {
		win.MovePrint(startY+row, startX, line)
	}

}
