package state

import (
	"gogue/model/entity"
	"gogue/model/signal"
	"gogue/presentation/action"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Game struct {
	player entity.Player
	view   *tview.Box
	signal signal.Type
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

	// @todo - выделить отрисовку в отдельный факл в view/cli
	box := tview.NewBox().SetBorder(true).SetTitle("Game")

	game := Game{
		player: player,
		view:   box,
	}

	game.view.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// -2: учёт рамки
		fw, fh := width-2, height-2
		game.drawField(screen, x+1, y+1, fw, fh)
		return x, y, width, height
	})

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
	// Нет lastField, всё строится на лету
	sig := g.signal
	g.signal = signal.NoSignal
	return sig
}

func (g *Game) Primitive() tview.Primitive {
	return g.view
}

func (g *Game) drawField(screen tcell.Screen, ox, oy, w, h int) {
	field := g.makeField(w, h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ch := ' '
			if field[y][x] == 1 {
				ch = '🦸'
			}
			screen.SetContent(ox+x, oy+y, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}
}

func (g *Game) makeField(w, h int) [][]int {
	field := make([][]int, h)
	for y := range field {
		field[y] = make([]int, w)
	}
	px := g.player.Character.Shape.Point.X
	py := g.player.Character.Shape.Point.Y
	if py >= 0 && py < h && px >= 0 && px < w {
		field[py][px] = 1
	}
	return field
}
