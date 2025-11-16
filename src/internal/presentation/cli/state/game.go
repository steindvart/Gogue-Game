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

const (
	HeightLevel = 30
	WidthLevel  = 90
)

type Symbol int

const (
	symPlayer Symbol = iota + 1
	symWall
	symPortal
	symPassage
	symDoor
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

	err := level.GenerateLevel(primitives.Size2D[uint]{Height: HeightLevel, Width: WidthLevel})
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

	movement := map[action.Type]primitives.Point2D[int]{
		action.MoveUp:               {X: 0, Y: -1},
		action.MoveDown:             {X: 0, Y: 1},
		action.MoveLeft:             {X: -1, Y: 0},
		action.MoveRight:            {X: 1, Y: 0},
		action.MoveLeftUpperCorner:  {X: -1, Y: -1},
		action.MoveRightUpperCorner: {X: 1, Y: -1},
		action.MoveLefLowerCorner:   {X: -1, Y: 1},
		action.MoveRightLowerCorner: {X: 1, Y: 1},
	}
	game.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch game.eventToAction(event) {
		case action.MoveUp,
			action.MoveDown,
			action.MoveLeft,
			action.MoveRight,
			action.MoveLeftUpperCorner,
			action.MoveRightUpperCorner,
			action.MoveLefLowerCorner,
			action.MoveRightLowerCorner:
			oldPosition := game.moveAndGetOldPosition(movement[game.eventToAction(event)])
			if game.checkCollision() {
				game.player.Character.Shape.Point = oldPosition
			}
		case action.Exit:
			game.signal = signals.Stop
		default:
			return event
		}
		return nil
	})

	return &game, nil
}

func (g *Game) moveAndGetOldPosition(changingPosition primitives.Point2D[int]) primitives.Point2D[int] {
	oldPosition := g.player.Character.Shape.Point
	g.player.Character.Shape.Move(changingPosition)
	return oldPosition
}

func (g *Game) checkCollision() bool {
	playerPosition := g.player.Character.Shape.Point

	// Простая проверка границ
	if playerPosition.X < 0 || playerPosition.Y < 0 ||
		playerPosition.X >= int(WidthLevel) || playerPosition.Y >= int(HeightLevel) {
		return true
	}

	// Проверка стен на карте
	for _, room := range g.level.Rooms {
		if g.isPlayerInWall(playerPosition, room) {
			return true
		}
	}

	return false
}

func (g *Game) isPlayerInWall(pos primitives.Point2D[int], room world.Room) bool {
	roomX := room.Shape.Point.X
	roomY := room.Shape.Point.Y
	roomWidth := int(room.Shape.Size.Width)
	roomHeight := int(room.Shape.Size.Height)

	for _, door := range room.Doors {
		if pos == door {
			return false
		}
	}

	// Верхняя или нижняя стена
	if pos.Y == roomY || pos.Y == roomY+roomHeight {
		return pos.X >= roomX && pos.X <= roomX+roomWidth
	}
	// Левая или правая стена
	if pos.X == roomX || pos.X == roomX+roomWidth {
		return pos.Y >= roomY && pos.Y <= roomY+roomHeight
	}

	return false
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
	field := g.makeField(w, h) // WidthLevel, HeightLevel
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ch := ' '
			if field[y][x] == int(symPlayer) {
				ch = '8' // '🦸'
			} else if field[y][x] == int(symWall) {
				ch = '⚀'
			} else if field[y][x] == int(symPortal) {
				ch = '0'
			} else if field[y][x] == int(symPassage) {
				ch = '*'
			} else if field[y][x] == int(symDoor) {
				ch = 'П'
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
		field[room.Shape.Point.Y][column] = int(symWall)
		field[room.Shape.Point.Y+int(room.Shape.Size.Height)][column] = int(symWall)
	}

	for row := room.Shape.Point.Y; row <= room.Shape.Point.Y+int(room.Shape.Size.Height); row++ {
		field[row][room.Shape.Point.X] = int(symWall)
		field[row][room.Shape.Point.X+int(room.Shape.Size.Width)] = int(symWall)
	}

	field[finishPortal.Point.Y][finishPortal.Point.X] = int(symPortal)
}

func (g *Game) drawPassage(passage world.Passage, field [][]int) {
	for i := 0; i < len(passage.Passage); i++ {
		field[passage.Passage[i].Y][passage.Passage[i].X] = int(symPassage)
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = int(symDoor)
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = int(symDoor)
}
