package state

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/model/signals"
	"gogue/internal/model/world"
	"gogue/internal/presentation/action"
	"math/rand"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Game struct {
	player *entities.Player
	level  *world.Level
	view   *tview.Box
	signal signals.Type
}

func NewGame() (*Game, error) {
	// @todo - выделить отрисовку в отдельный файл в view/cli
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := world.NewLevel(source)

	err := level.GenerateLevel(primitives.Size2D[uint]{Height: 30, Width: 90})
	if err != nil {
		return nil, err
	}

	playerStartPoint, err := level.GetStartPositionForPlayer()
	if err != nil {
		return nil, err
	}
	player := entities.NewPlayer(primitives.Box{
		Point: *playerStartPoint,
		Size:  primitives.Size2D[uint]{Height: 1, Width: 1},
	})

	box := tview.NewBox().SetBorder(true).SetTitle("Game")

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
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: 0, Y: -1})
			return nil
		case action.MoveDown:
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: 0, Y: 1})
			return nil
		case action.MoveLeft:
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: -1, Y: 0})
			return nil
		case action.MoveRight:
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: 1, Y: 0})
			return nil
		case action.MoveLeftUpperCorner:
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: -1, Y: -1})
			return nil
		case action.MoveRightUpperCorner:
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: 1, Y: -1})
			return nil
		case action.MoveLefLowerCorner:
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: -1, Y: 1})
			return nil
		case action.MoveRightLowerCorner:
			game.player.Character.Shape.Move(primitives.Point2D[int]{X: 1, Y: 1})
			return nil
		case action.Exit:
			game.signal = signals.Stop
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

func (g *Game) Update(float64) signals.Type {
	// Нет lastField, всё строится на лету
	sig := g.signal
	g.signal = signals.NoSignal
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
			} else if field[y][x] == 5 {
				ch = '*'
			} else if field[y][x] == 6 {
				ch = '['
			} else if field[y][x] == 7 {
				ch = ']'
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
		g.drawRoom(room, g.level.FinishPortal, field)
	}

	for _, passages := range g.level.Passages {
		g.drawPassage(passages, field)
	}

	px := g.player.Character.Shape.Point.X
	py := g.player.Character.Shape.Point.Y
	if py >= 0 && py < h && px >= 0 && px < w {
		field[py][px] = 1
	}

	return field
}

func (g *Game) drawRoom(room world.Room, finishPortal primitives.Box, field [][]int) {
	for column := room.Shape.Point.X; column < room.Shape.Point.X+int(room.Shape.Size.Width); column++ {
		field[room.Shape.Point.Y][column] = 2
		field[room.Shape.Point.Y+int(room.Shape.Size.Height)][column] = 2
	}

	for row := room.Shape.Point.Y; row < room.Shape.Point.Y+int(room.Shape.Size.Height); row++ {
		field[row][room.Shape.Point.X] = 3
		field[row][room.Shape.Point.X+int(room.Shape.Size.Width)] = 3
	}

	field[finishPortal.Point.Y][finishPortal.Point.X] = 4
}

func (g *Game) drawPassage(passage world.Passage, field [][]int) {
	for i := 0; i < len(passage.Passage); i++ {
		field[passage.Passage[i].Y][passage.Passage[i].X] = 5
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = 6
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = 7
}
