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
	HeightLevel = 30
	WidthLevel  = 90
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
			if game.checkCollision(oldPosition) {
				game.player.Character.Shape.Point = oldPosition
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

func (g *Game) moveAndGetOldPosition(changingPosition primitives.Point2D[int]) primitives.Point2D[int] {
	oldPosition := g.player.Character.Shape.Point
	g.player.Character.Shape.Move(changingPosition)
	return oldPosition
}

func (g *Game) checkCollision(oldPosition primitives.Point2D[int]) bool {
	playerPosition := g.player.Character.Shape.Point

	if playerPosition.X < 0 || playerPosition.Y < 0 {
		return true
	}
	if playerPosition.X >= WidthLevel || playerPosition.Y >= HeightLevel {
		return true
	}
	if g.checkCollisionInRoom() {
		return true
	}
	if g.checkCollisionWithEnemy() {
		return true
	}
	if inPassage, hasCollision := g.checkCollisionInPassage(oldPosition); inPassage {
		if !hasCollision {
			return true
		}
	}

	return false
}

func (g *Game) checkCollisionInRoom() bool {
	playerPosition := g.player.Character.Shape.Point

	for _, room := range g.level.Rooms {
		if isPlayerInWall(playerPosition, room) {
			return true
		}
	}

	return false
}

func isPlayerInWall(pos primitives.Point2D[int], room world.Room) bool {
	roomX := room.Shape.Point.X
	roomY := room.Shape.Point.Y
	roomWidth := int(room.Shape.Size.Width)
	roomHeight := int(room.Shape.Size.Height)

	for _, door := range room.Doors {
		if pos == door {
			return false
		}
	}

	if pos.Y == roomY || pos.Y == roomY+roomHeight {
		return pos.X >= roomX && pos.X <= roomX+roomWidth
	}
	if pos.X == roomX || pos.X == roomX+roomWidth {
		return pos.Y >= roomY && pos.Y <= roomY+roomHeight
	}

	return false
}

func (g *Game) checkCollisionInPassage(oldPosition primitives.Point2D[int]) (inPassage bool, hasCollision bool) {
	playerPosition := g.player.Character.Shape.Point

	for _, passages := range g.level.Passages {
		for _, passagePoint := range passages.Passage {
			if playerPosition == passagePoint {
				return true, true
			}
		}
		if playerPosition == passages.DoorOne || playerPosition == passages.DoorTwo {
			return true, true
		}
	}

	for _, passages := range g.level.Passages {
		for _, passagePoint := range passages.Passage {
			if oldPosition == passagePoint {
				return true, false
			}
		}
	}

	return false, false
}

func (g *Game) checkCollisionWithEnemy() bool {
	playerPosition := g.player.Character.Shape.Point

	for _, room := range g.level.Rooms {
		for _, enemy := range room.Enemies {
			if playerPosition == enemy.Character.Shape.Point {
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
		{
			Type: entities.EnemyType(entities.EnemyTypeZombie),
			Character: entities.Character{
				Shape: primitives.Box{
					Point: primitives.Point2D[int]{X: 6, Y: 6},
				},
			},
		},
		{
			Type: entities.EnemyType(entities.EnemyTypeVampire),
			Character: entities.Character{
				Shape: primitives.Box{
					Point: primitives.Point2D[int]{X: 6, Y: 7},
				},
			},
		},
		{
			Type: entities.EnemyType(entities.EnemyTypeGhost),
			Character: entities.Character{
				Shape: primitives.Box{
					Point: primitives.Point2D[int]{X: 6, Y: 8},
				},
			},
		},
		{
			Type: entities.EnemyType(entities.EnemyTypeOgre),
			Character: entities.Character{
				Shape: primitives.Box{
					Point: primitives.Point2D[int]{X: 6, Y: 9},
				},
			},
		},
		{
			Type: entities.EnemyType(entities.EnemyTypeSnakeMage),
			Character: entities.Character{
				Shape: primitives.Box{
					Point: primitives.Point2D[int]{X: 6, Y: 10},
				},
			},
		},
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
	for col := room.Shape.Point.X; col <= room.Shape.Point.X+width; col++ {
		field[room.Shape.Point.Y][col] = common.EntityTypeWall
		field[room.Shape.Point.Y+height][col] = common.EntityTypeWall
	}

	for row := room.Shape.Point.Y; row < room.Shape.Point.Y+height; row++ {
		field[row][room.Shape.Point.X] = common.EntityTypeWall
		field[row][room.Shape.Point.X+width] = common.EntityTypeWall
	}

	field[finishPortal.Point.Y][finishPortal.Point.X] = common.EntityTypePortal
}

func (g *Game) putPassage(passage world.Passage, field [][]common.EntityType) {
	for i := 0; i < len(passage.Passage); i++ {
		field[passage.Passage[i].Y][passage.Passage[i].X] = common.EntityTypePassage
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = common.EntityTypeDoor
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = common.EntityTypeDoor
}
