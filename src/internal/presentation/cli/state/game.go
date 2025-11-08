package state

import (
	"gogue/internal/model/entity"
	"gogue/internal/model/primitive"
	"gogue/internal/model/signal"
	"gogue/internal/model/world"
	"gogue/internal/presentation/action"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Game struct {
	player entity.Player
	level  world.Level
	view   *tview.Box
	signal signal.Type
}

func NewGame() (*Game, error) {
	player := entity.Player{
		Character: entity.Character{
			Shape: primitive.Box{
				Point: primitive.Point2D[int]{X: 5, Y: 5},
				Size:  primitive.Size2D[uint]{Height: 1, Width: 1},
			},
			Attributes: primitive.Attributes{
				Health:    100,
				MaxHealth: 100,
				Strength:  10,
				Agility:   5,
			},
		},
		Backpack: nil,
		Weapon:   nil,
	}

	// @todo - выделить отрисовку в отдельный файл в view/cli
	box := tview.NewBox().SetBorder(true).SetTitle("Game")

	level := world.Level{}
	// @todo - обработать ошибку
	err := level.GenerateNineRooms(primitive.Size2D[uint]{Height: 30, Width: 90})
	if err != nil {
		return nil, err
	}

	game := Game{
		player: player,
		level:  level,
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
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: 0, Y: -1})
			return nil
		case action.MoveDown:
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: 0, Y: 1})
			return nil
		case action.MoveLeft:
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: -1, Y: 0})
			return nil
		case action.MoveRight:
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: 1, Y: 0})
			return nil
		case action.MoveLeftUpperCorner:
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: -1, Y: -1})
			return nil
		case action.MoveRightUpperCorner:
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: 1, Y: -1})
			return nil
		case action.MoveLefLowerCorner:
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: -1, Y: 1})
			return nil
		case action.MoveRightLowerCorner:
			game.player.Character.Shape.Move(primitive.Point2D[int]{X: 1, Y: 1})
			return nil
		case action.Exit:
			game.signal = signal.Stop
			return nil
		default:
			return event
		}
	})

	return &game, nil
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

	ch := unicode.ToLower(event.Rune())
	switch ch {
	case 'w', 'ц':
		return action.MoveUp
	case 's', 'ы':
		return action.MoveDown
	case 'a', 'ф':
		return action.MoveLeft
	case 'd', 'в':
		return action.MoveRight
	case 'y', 'н':
		return action.MoveLeftUpperCorner
	case 'u', 'г':
		return action.MoveRightUpperCorner
	case 'b', 'и':
		return action.MoveLefLowerCorner
	case 'n', 'т':
		return action.MoveRightLowerCorner
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
			} else if field[y][x] == 2 {
				ch = '—'
			} else if field[y][x] == 3 {
				ch = '|'
			} else if field[y][x] == 4 {
				ch = 'O'
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

	for _, room := range g.level.Rooms {
		g.drawRoom(room, g.level, field)
	}

	px := g.player.Character.Shape.Point.X
	py := g.player.Character.Shape.Point.Y
	if py >= 0 && py < h && px >= 0 && px < w {
		field[py][px] = 1
	}

	return field
}

func (g *Game) drawRoom(room world.Room, level world.Level, field [][]int) {
	for column := room.Shape.Point.X; column < room.Shape.Point.X+int(room.Shape.Size.Width); column++ {
		field[room.Shape.Point.Y][column] = 2
		field[room.Shape.Point.Y+int(room.Shape.Size.Height)][column] = 2
	}

	for row := room.Shape.Point.Y; row < room.Shape.Point.Y+int(room.Shape.Size.Height); row++ {
		field[row][room.Shape.Point.X] = 3
		field[row][room.Shape.Point.X+int(room.Shape.Size.Width)] = 3
	}

	field[level.End.Point.Y][level.End.Point.X] = 4
}
