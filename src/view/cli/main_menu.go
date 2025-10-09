package cli

import (
	"errors"
	"fmt"
	"time"

	gc "github.com/rthornton128/goncurses"
)

const (
	MENU_HEIGHT = 35
	MENU_WIDTH  = 70
)

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

var hallArt = []string{
	" _____________________________________________",
	"|.'',                                     ,''.|",
	"|.'.'',                                 ,''.'.|",
	"|.'.'.'',                             ,''.'.'.|",
	"|.'.'.'.'',                         ,''.'.'.'.|",
	"|.'.'.'.'.|                         |.'.'.'.'.|",
	"|.'.'.'.'.|===;                 ;===|.'.'.'.'.|",
	"|.'.'.'.'.|:::|',             ,'|:::|.'.'.'.'.|",
	"|.'.'.'.'.|---|'.|, _______ ,|.'|---|.'.'.'.'.|",
	"|.'.'.'.'.|:::|'.|'|???????|'|.'|:::|.'.'.'.'.|",
	"|,',',',',|---|',|'|???????|'|,'|---|,',',',',|",
	"|.'.'.'.'.|:::|'.|'|???????|'|.'|:::|.'.'.'.'.|",
	"|.'.'.'.'.|---|','   /%%%\\   ','|---|.'.'.'.'.|",
	"|.'.'.'.'.|===:'    /%%%%%\\    ':===|.'.'.'.'.|",
	"|.'.'.'.'.|%%%%%%%%%%%%%%%%%%%%%%%%%|.'.'.'.'.|",
	"|.'.'.'.','       /%%%%%%%%%\\       ','.'.'.'.|",
	"|.'.'.','        /%%%%%%%%%%%\\        ','.'.'.|",
	"|.'.','         /%%%%%%%%%%%%%\\         ','.'.|",
	"|.','          /%%%%%%%%%%%%%%%\\          ','.|",
	"|;____________/%%%%%%%%%%%%%%%%%\\____________;|",
}

type MainMenu struct {
	window *gc.Window
	colors []int
}

var (
	errMenuAttrOn        = errors.New("cannot set menu attribute on")
	errMenuAttrOff       = errors.New("cannot set menu attribute off")
	errMenuInitColorPair = errors.New("cannot init color pair")
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
		colors, err = create256RainbowPairs()
		if err != nil {
			return nil, err
		}
	}

	return &MainMenu{window: win, colors: colors}, nil
}

func create256RainbowPairs() ([]int, error) {
	var rainbow256 = []int{196, 202, 208, 214, 220, 226, 190, 154, 118, 82, 46, 47, 48, 49, 51, 39, 27, 21, 57, 93, 129, 165, 201, 200}

	for i, color := range rainbow256 {
		err := gc.InitPair(int16(i+1), int16(color), gc.C_BLACK)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMenuInitColorPair, err)
		}
	}
	return rainbow256, nil
}

func (m *MainMenu) Render(options []string, active int) error {
	m.RenderTitle()

	err := m.RenderOptions(options, active)
	if err != nil {
		return err
	}

	// @todo - пока убрал, т.к. кажется не совсем уместным, ломает минималистичный стиль
	// Но если нравится - можем оставить. Включите, посмотрите, как с этим будет смотреться.
	// Также можем раскрасить как-нибудь.
	// m.RenderHall()

	m.window.Refresh()
	return nil
}

func (m *MainMenu) RenderHall() {
	maxY, maxX := m.window.MaxYX()
	artHeight := len(hallArt)
	artWidth := len(hallArt[0])

	startY := maxY - artHeight - 1 // maxY - рисуем внизу окна
	startX := (maxX - artWidth) / 2

	for i, line := range hallArt {
		m.window.MovePrint(startY+i, startX, line)
	}
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

func (m *MainMenu) RenderTitle() error {
	if gc.HasColors() {
		// Для плавности: используем текущее время как frame
		frame := int(time.Now().UnixNano() / 85000000) // ~13 кадров в сек
		err := drawGradientTitle(m.window, frame, m.colors)
		if err != nil {
			return err
		}
	} else {
		drawRawTitle(m.window)
	}

	return nil
}

func drawGradientTitle(win *gc.Window, frame int, colors []int) error {
	_, maxX := win.MaxYX()
	titleWidth := len(titleArt[0])
	startY := 1                       // чуть ниже начала окна
	startX := (maxX - titleWidth) / 2 // по центру с учётом ширины

	for row, line := range titleArt {
		for col, ch := range line {
			// Диагональное переливание: сдвиг также по row
			colorIdx := (col + row + frame) % len(colors)

			err := win.AttrOn(gc.ColorPair(int16(colorIdx + 1)))
			if err != nil {
				return fmt.Errorf("%w: %v", errMenuAttrOn, err)
			}

			win.MovePrint(startY+row, startX+col, string(ch))

			err = win.AttrOff(gc.ColorPair(int16(colorIdx + 1)))
			if err != nil {
				return fmt.Errorf("%w: %v", errMenuAttrOn, err)
			}
		}
	}

	return nil
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
