package state

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
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
	MapHeight = 30
	MapWidth  = 90
)

type Game struct {
	level  *world.Level
	view   *viewcli.Game
	signal signals.Type
}

func NewGame() (*Game, error) {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := world.NewLevel(source, primitives.Size2D[uint]{Height: MapHeight, Width: MapWidth})

	level.Generate()

	game := Game{
		level: level,
		view:  viewcli.NewGame(),
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
			level.MovePlayer(movementRegistry[game.eventToAction(event)])
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

func (g *Game) makeField(w, h int) [][]common.GameEntityType {
	field := make([][]common.GameEntityType, h)
	for y := range field {
		field[y] = make([]common.GameEntityType, w)
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
		g.putItems(room, field)
	}

	for _, passages := range g.level.Passages {
		g.putPassage(passages, field)
	}

	// for _, room := range g.level.Rooms {
	// 	g.putEnemies(room, field)
	// }

	px := g.level.Player.Character.Shape.Point.X
	py := g.level.Player.Character.Shape.Point.Y
	if py >= 0 && py < h && px >= 0 && px < w {
		field[py][px] = common.EntityTypePlayer
	}

	return field
}

//func (g *Game) putEnemies(room world.Room, field [][]common.GameEntityType) {
//	var et common.GameEntityType
//
//	for _, e := range room.Enemies {
//		switch e.Type {
//		case entities.EnemyTypeZombie:
//			et = common.EntityTypeZombie
//		case entities.EnemyTypeVampire:
//			et = common.EntityTypeVampire
//		case entities.EnemyTypeGhost:
//			et = common.EntityTypeGhost
//		case entities.EnemyTypeOgre:
//			et = common.EntityTypeOgre
//		case entities.EnemyTypeSnakeMage:
//			et = common.EntityTypeSnakeMage
//		}
//
//		ex := e.Character.Shape.Point.X
//		ey := e.Character.Shape.Point.Y
//
//		h := len(field)
//		w := len(field[0])
//		if ey >= 0 && ey < h && ex >= 0 && ex < w {
//			field[ey][ex] = et
//		}
//	}
//}

// Тут можно класть только lvl, так как room я получаю из него же шагом выше, а могу и тут
func (g *Game) putRoom(room world.Room, finishPortal primitives.Box, field [][]common.GameEntityType) {
	width := int(room.Shape.Size.Width)
	height := int(room.Shape.Size.Height)

	startX := room.Shape.Point.X
	endX := startX + width - 1
	startY := room.Shape.Point.Y
	endY := startY + height - 1

	for col := startX; col <= endX; col++ {
		field[startY][col] = common.WorldTypeWall
		field[endY][col] = common.WorldTypeWall
	}

	for row := startY; row <= endY; row++ {
		field[row][startX] = common.WorldTypeWall
		field[row][endX] = common.WorldTypeWall
	}

	field[finishPortal.Point.Y][finishPortal.Point.X] = common.WorldTypePortal
}

func (g *Game) putPassage(passage world.Passage, field [][]common.GameEntityType) {
	for i := 0; i < len(passage.Way); i++ {
		field[passage.Way[i].Y][passage.Way[i].X] = common.WorldTypePassage
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = common.WorldTypeDoor
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = common.WorldTypeDoor
}

// @todo - оставить ли это всё кучей?
func (g *Game) putItems(room world.Room, field [][]common.GameEntityType) {
	var fd common.GameEntityType
	for _, food := range room.Foods {
		switch food.Type {
		case items.FoodTypePotatoes:
			fd = common.FoodTypePotatoes
		case items.FoodTypeBread:
			fd = common.FoodTypeBread
		case items.FoodTypeMeat:
			fd = common.FoodTypeMeat
		case items.FoodTypeMistery:
			fd = common.FoodTypeMistery
		case items.FoodTypeBeer:
			fd = common.FoodTypeBeer
		}
		field[food.Item.Shape.Point.Y][food.Item.Shape.Point.X] = fd
	}

	for _, elixir := range room.Elixirs {
		switch elixir.Type {
		case items.ElixirTypeStrength:
			fd = common.ElixirTypeStrength
		case items.ElixirTypeAgility:
			fd = common.ElixirTypeAgility
		case items.ElixirTypeDwarfism:
			fd = common.ElixirTypeDwarfism
		case items.ElixirTypeGiantism:
			fd = common.ElixirTypeGiantism
		case items.ElixirTypeMystery:
			fd = common.ElixirTypeMystery
		}
		field[elixir.Item.Shape.Point.Y][elixir.Item.Shape.Point.X] = fd
	}

	for _, scroll := range room.Scrolls {
		switch scroll.Type {
		case items.ScrollTypeStrength:
			fd = common.ScrollTypeStrength
		case items.ScrollTypeAgility:
			fd = common.ScrollTypeAgility
		case items.ScrollTypeUltimate:
			fd = common.ScrollTypeUltimate
		case items.ScrollTypeMaxHealth:
			fd = common.ScrollTypeMaxHealth
		case items.ScrollTypeMystery:
			fd = common.ScrollTypeMystery
		}
		field[scroll.Item.Shape.Point.Y][scroll.Item.Shape.Point.X] = fd
	}

	for _, weapon := range room.Weapons {
		fd = common.Weapon
		field[weapon.Item.Shape.Point.Y][weapon.Item.Shape.Point.X] = fd
	}
}
