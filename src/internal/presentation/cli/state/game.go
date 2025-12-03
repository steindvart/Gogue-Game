package state

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/model/signals"
	"gogue/internal/model/world"
	"gogue/internal/presentation/action"
	viewcli "gogue/internal/view/cli"
	"math/rand"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	LevelHeight = 30
	LevelWidth  = 90
)

type Game struct {
	player *entities.Player
	level  *world.Level
	view   *viewcli.Game
	signal signals.Type
}

func NewGame() (*Game, error) {
	// @todo - выделить отрисовку в отдельный файл в view/cli
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := world.NewLevel(source)

	err := level.GenerateLevel(primitives.Size2D[uint]{Height: LevelHeight, Width: LevelWidth})
	if err != nil {
		return nil, err
	}

	startPlayerPos, err := level.GenerateStartPlayerPosition()
	if err != nil {
		return nil, err
	}
	player := entities.NewPlayer(&primitives.Box{
		Point: *startPlayerPos,
		Size:  primitives.Size2D[uint]{Height: 1, Width: 1},
	})

	gameView := viewcli.NewGame()

	game := Game{
		player: player,
		level:  level,
		view:   gameView,
	}

	game.view.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// -2: учёт рамки
		fw, fh := width-2, height-2
		field := game.makeField(fw, fh)
		game.view.SetFieldToScreen(screen, field, x+1, y+1)
		return x, y, width, height
	})

	movementRegistry := map[action.Type]primitives.Point2D[int]{
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

			oldPlayerPos := game.player.GetPosition()
			game.player.Move(movementRegistry[game.eventToAction(event)])

			if game.checkCollision(game.player.GetPosition()) {
				game.player.SetPosition(oldPlayerPos)
			}

		case action.Exit:
			game.signal = signals.Stop
			return nil
		default:
			return event
		}
		return nil
	})

	return &game, nil
}

func (g *Game) checkCollision(pos primitives.Point2D[int]) bool {
	if g.checkCollisionWithFieldBorders(pos) {
		return true
	}
	if g.checkCollisionWithRoomsWall(pos) {
		return true
	}
	if g.checkCollisionWithEnemy(pos) {
		return true
	}

	if !isInSomeRoom(pos, g.level.Rooms) && !g.checkCollisionWithPassages(pos) {
		return true
	}

	return false
}

func (g *Game) checkCollisionWithFieldBorders(pos primitives.Point2D[int]) bool {
	if pos.X < 0 || pos.Y < 0 {
		return true
	}
	if pos.X >= LevelWidth || pos.Y >= LevelHeight {
		return true
	}

	return false
}

func (g *Game) checkCollisionWithRoomsWall(pos primitives.Point2D[int]) bool {
	for _, room := range g.level.Rooms {
		if isInRoom(pos, room) && checkCollisionWithRoomWall(pos, room) {
			return true
		}
	}

	return false
}

func isInSomeRoom(pos primitives.Point2D[int], rooms []world.Room) bool {
	for _, room := range rooms {
		if isInRoom(pos, room) {
			return true
		}
	}

	return false
}

func isInRoom(pos primitives.Point2D[int], room world.Room) bool {
	leftEndX := room.Shape.Point.X + 1
	rightEndX := leftEndX + int(room.Shape.Size.Width) - 3
	topEndY := room.Shape.Point.Y + 1
	downEndY := topEndY + int(room.Shape.Size.Height) - 3

	return (pos.X >= leftEndX && pos.X <= rightEndX) &&
		(pos.Y >= topEndY && pos.Y <= downEndY)
}

func checkCollisionWithRoomWall(pos primitives.Point2D[int], room world.Room) bool {
	// Двери явлюятся частью комнаты и её стен, но через них можно ходить
	if checkCollisionWithDoors(pos, room.Doors) {
		return false
	}

	leftEndX := room.Shape.Point.X
	rightEndX := leftEndX + int(room.Shape.Size.Width)
	topEndY := room.Shape.Point.Y
	downEndY := topEndY + int(room.Shape.Size.Height)

	if (pos.X == leftEndX || pos.X == rightEndX) || (pos.Y == topEndY || pos.Y == downEndY) {
		return true
	}

	return false
}

func checkCollisionWithDoors(pos primitives.Point2D[int], doors []primitives.Point2D[int]) bool {
	for _, door := range doors {
		if pos == door {
			return true
		}
	}

	return false
}

func (g *Game) checkCollisionWithPassages(newPos primitives.Point2D[int]) bool {
	for _, passage := range g.level.Passages {
		if isInPassage(newPos, passage) {
			return true
		}
	}

	return false
}

func isInPassage(pos primitives.Point2D[int], passage world.Passage) bool {
	// Двери являются как частью комнаты, так и частью прохода
	if pos == passage.DoorOne || pos == passage.DoorTwo {
		return true
	}

	for _, wayPoint := range passage.Way {
		if pos == wayPoint {
			return true
		}
	}

	return false
}

func (g *Game) checkCollisionWithEnemy(pos primitives.Point2D[int]) bool {
	for _, room := range g.level.Rooms {
		for _, enemy := range room.Enemies {
			if pos == enemy.GetPosition() {
				return true
			}
		}
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

func (g *Game) makeField(w, h int) [][]common.EntityType {
	field := make([][]common.EntityType, h)
	for y := range field {
		field[y] = make([]common.EntityType, w)
	}

	// Добавлено в качестве примера, потом надо будет убрать
	g.level.Rooms[0].Enemies = []entities.Enemy{
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeZombie),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 6},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeVampire),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 7},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeGhost),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 8},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeOgre),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 9},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeSnakeMage),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 10},
		//		},
		//	},
		//},
	}

	for _, room := range g.level.Rooms {
		g.putRoom(room, g.level.FinishPortal, field)
	}

	for _, passages := range g.level.Passages {
		g.putPassage(passages, field)
	}

	for _, room := range g.level.Rooms {
		g.putEnemies(room, field)
	}

	px := g.player.Character.Shape.Point.X
	py := g.player.Character.Shape.Point.Y
	if py >= 0 && py < h && px >= 0 && px < w {
		field[py][px] = common.EntityTypePlayer
	}

	return field
}

func (g *Game) putEnemies(room world.Room, field [][]common.EntityType) {
	var et common.EntityType

	for _, e := range room.Enemies {
		switch e.Type {
		case entities.EnemyTypeZombie:
			et = common.EntityTypeZombie
		case entities.EnemyTypeVampire:
			et = common.EntityTypeVampire
		case entities.EnemyTypeGhost:
			et = common.EntityTypeGhost
		case entities.EnemyTypeOgre:
			et = common.EntityTypeOgre
		case entities.EnemyTypeSnakeMage:
			et = common.EntityTypeSnakeMage
		}

		ex := e.Character.Shape.Point.X
		ey := e.Character.Shape.Point.Y

		h := len(field)
		w := len(field[0])
		if ey >= 0 && ey < h && ex >= 0 && ex < w {
			field[ey][ex] = et
		}
	}
}

// Тут можно класть только lvl, так как room я получаю из него же шагом выше, а могу и тут
func (g *Game) putRoom(room world.Room, finishPortal primitives.Box, field [][]common.EntityType) {
	width := int(room.Shape.Size.Width)
	height := int(room.Shape.Size.Height)

	startX := room.Shape.Point.X
	endX := startX + width - 1
	startY := room.Shape.Point.Y
	endY := startY + height - 1

	for col := startX; col <= endX; col++ {
		field[startY][col] = common.EntityTypeWall
		field[endY][col] = common.EntityTypeWall
	}

	for row := startY; row <= endY; row++ {
		field[row][startX] = common.EntityTypeWall
		field[row][endX] = common.EntityTypeWall
	}

	field[finishPortal.Point.Y][finishPortal.Point.X] = common.EntityTypePortal
}

func (g *Game) putPassage(passage world.Passage, field [][]common.EntityType) {
	for i := 0; i < len(passage.Way); i++ {
		field[passage.Way[i].Y][passage.Way[i].X] = common.EntityTypePassage
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = common.EntityTypeDoor
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = common.EntityTypeDoor
}
