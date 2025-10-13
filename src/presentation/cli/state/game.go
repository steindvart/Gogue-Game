package state

import (
	"gogue/model/entity"
	"gogue/model/signal"
	"gogue/view/action"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	FieldHeight = 20
	FieldWidth  = 40
)

type Game struct {
	player      entity.Player
	fieldWidth  int
	fieldHeight int
	view        *tview.Box
	lastField   [][]int
	signal      signal.Type
}

func NewGame() *Game {
	player := entity.Player{
		Character: entity.Character{
			Shape: entity.Box{
				Point: entity.Point2D[int]{X: 5, Y: 5},
				Size:  entity.Size2D[uint]{Height: 1, Width: 1},
			},
			Health:    100,
			MaxHealth: 100,
			Strength:  10,
			Agility:   5,
		},
		Backpack: nil,
		Weapon:   nil,
	}

	box := tview.NewBox().SetBorder(true).SetTitle("Game")

	game := Game{
		player:      player,
		fieldWidth:  FieldWidth,
		fieldHeight: FieldHeight,
		view:        box,
	}

	game.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch game.eventToAction(event) {
		case action.MoveUp:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: 0, Y: -1})
			return nil
		case action.MoveDown:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: 0, Y: 1})
			return nil
		case action.MoveLeft:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: -1, Y: 0})
			return nil
		case action.MoveRight:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: 1, Y: 0})
			return nil
		case action.Exit:
			game.signal = signal.Stop
			return nil
		default:
			return event
		}
	})

	game.view.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		game.view.Draw(screen)
		game.drawField(screen, x+1, y+1, width-2, height-2)
		return x, y, width, height
	})

	return &game
}

func (g *Game) eventToAction(event *tcell.EventKey) action.Type {
	switch event.Key() {
	case tcell.KeyUp:
		return action.MoveUp
	case tcell.KeyDown:
		return action.MoveDown
	case tcell.KeyLeft:
		return action.MoveLeft
	case tcell.KeyRight:
		return action.MoveRight
	case tcell.KeyEsc:
		return action.Exit
	}

	switch event.Rune() {
	case 'w', 'W', 'ц', 'Ц':
		return action.MoveUp
	case 's', 'S', 'ы', 'Ы':
		return action.MoveDown
	case 'a', 'A', 'ф', 'Ф':
		return action.MoveLeft
	case 'd', 'D', 'в', 'В':
		return action.MoveRight
	}

	return action.NoAction
}

func (g *Game) Update(float64) signal.Type {
	g.lastField = g.makeField()

	sig := g.signal
	g.signal = signal.NoSignal // сброс сигнала после чтения
	return sig
}

func (g *Game) Primitive() tview.Primitive {
	return g.view
}

func (g *Game) drawField(screen tcell.Screen, ox, oy, w, h int) {
	field := g.lastField
	if field == nil {
		field = g.makeField()
	}

	for y := 0; y < g.fieldHeight && y < h; y++ {
		for x := 0; x < g.fieldWidth && x < w; x++ {
			ch := ' '
			if field[y][x] == 1 {
				ch = '🦸'
			}
			screen.SetContent(ox+x, oy+y, ch, nil, tcell.StyleDefault)
		}
	}
}

func (g *Game) makeField() [][]int {
	field := make([][]int, g.fieldHeight)
	for y := range field {
		field[y] = make([]int, g.fieldWidth)
	}
	px := g.player.Character.Shape.Point.X
	py := g.player.Character.Shape.Point.Y
	if py >= 0 && py < g.fieldHeight && px >= 0 && px < g.fieldWidth {
		field[py][px] = 1
	}
	return field
}
