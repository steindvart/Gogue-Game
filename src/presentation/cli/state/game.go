package state

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gogue/model/entity"
	"gogue/model/signal"
	"gogue/presentation/action"
	"math/rand"
	"time"
	"unicode"
)

type Game struct {
	player entity.Player
	level  *entity.Level
	view   *tview.Box
	signal signal.Type
}

func NewGame() (*Game, error) {
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

	// @todo - выделить отрисовку в отдельный файл в view/cli
	box := tview.NewBox().SetBorder(true).SetTitle("Game")

	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := entity.NewLevel(source)
	err := level.GenerateNineRooms(entity.Size2D[uint]{Height: 30, Width: 90})
	if err != nil {
		return nil, err
	}

	err = level.GeneratePassages()
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
		case action.MoveLeftUpperCorner:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: -1, Y: -1})
			return nil
		case action.MoveRightUpperCorner:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: 1, Y: -1})
			return nil
		case action.MoveLefLowerCorner:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: -1, Y: 1})
			return nil
		case action.MoveRightLowerCorner:
			game.player.Character.Shape.Move(entity.Point2D[int]{X: 1, Y: 1})
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

// Тут можно класть только lvl, так как room я получаю из него же шагом выше, а могу и тут
func (g *Game) drawRoom(room entity.Room, finishPortal entity.Box, field [][]int) {
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

func (g *Game) drawPassage(passage entity.Passage, field [][]int) {
	for i := 0; i < len(passage.Passage); i++ {
		field[passage.Passage[i].Y][passage.Passage[i].X] = 5
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = 6
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = 7
}
