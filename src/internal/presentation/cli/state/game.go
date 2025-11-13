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

type EntityType int

const (
	EntityTypePlayer EntityType = iota + 1
	EntityTypeHorizontalWall
	EntityTypeVerticalWall
	EntityTypePortal
	EntityTypePassage
	EntityTypeDoorOne
	EntityTypeDoorTwo
	EntityTypeZombie
	EntityTypeVampire
	EntityTypeGhost
	EntityTypeOgre
	EntityTypeSnakeMage
)

type Game struct {
	player entities.Player
	level  *world.Level
	view   *tview.Box
	signal signals.Type
}

func NewGame() (*Game, error) {
	player := entities.Player{
		Character: entities.Character{
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 5, Y: 5},
				Size:  primitives.Size2D[uint]{Height: 1, Width: 1},
			},
			Attributes: primitives.Attributes{
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

	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := world.NewLevel(source)
	err := level.GenerateNineRooms(primitives.Size2D[uint]{Height: 30, Width: 90})
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
			if field[y][x] == EntityTypePlayer {
				ch = '🦸'
			} else if field[y][x] == EntityTypeHorizontalWall {
				ch = '—'
			} else if field[y][x] == EntityTypeVerticalWall {
				ch = '|'
			} else if field[y][x] == EntityTypePortal {
				ch = 'O'
			} else if field[y][x] == EntityTypePassage {
				ch = '*'
			} else if field[y][x] == EntityTypeDoorOne {
				ch = '['
			} else if field[y][x] == EntityTypeDoorTwo {
				ch = ']'
			} else if field[y][x] == EntityTypeZombie {
				ch = '🧟'
			} else if field[y][x] == EntityTypeVampire {
				ch = '🧛'
			} else if field[y][x] == EntityTypeGhost {
				ch = '👻'
			} else if field[y][x] == EntityTypeOgre {
				ch = '👹'
			} else if field[y][x] == EntityTypeSnakeMage {
				ch = '🐍'
			}
			screen.SetContent(ox+x, oy+y, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}
}

func (g *Game) makeField(w, h int) [][]EntityType {
	field := make([][]EntityType, h)
	for y := range field {
		field[y] = make([]EntityType, w)
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
		g.drawRoom(room, g.level.FinishPortal, field)
		g.drawEnemies(room, field)
	}

	for _, passages := range g.level.Passages {
		g.drawPassage(passages, field)
	}

	px := g.player.Character.Shape.Point.X
	py := g.player.Character.Shape.Point.Y
	if py >= 0 && py < h && px >= 0 && px < w {
		field[py][px] = EntityTypePlayer
	}

	return field
}

func (g *Game) drawEnemies(room world.Room, field [][]EntityType) {
	var et EntityType

	for _, e := range room.Enemies {
		switch e.Type {
		case entities.EnemyTypeZombie:
			et = EntityTypeZombie
		case entities.EnemyTypeVampire:
			et = EntityTypeVampire
		case entities.EnemyTypeGhost:
			et = EntityTypeGhost
		case entities.EnemyTypeOgre:
			et = EntityTypeOgre
		case entities.EnemyTypeSnakeMage:
			et = EntityTypeSnakeMage
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
func (g *Game) drawRoom(room world.Room, finishPortal primitives.Box, field [][]EntityType) {
	for column := room.Shape.Point.X; column < room.Shape.Point.X+int(room.Shape.Size.Width); column++ {
		field[room.Shape.Point.Y][column] = EntityTypeHorizontalWall
		field[room.Shape.Point.Y+int(room.Shape.Size.Height)][column] = EntityTypeHorizontalWall
	}

	for row := room.Shape.Point.Y; row < room.Shape.Point.Y+int(room.Shape.Size.Height); row++ {
		field[row][room.Shape.Point.X] = EntityTypeVerticalWall
		field[row][room.Shape.Point.X+int(room.Shape.Size.Width)] = EntityTypeVerticalWall
	}

	field[finishPortal.Point.Y][finishPortal.Point.X] = EntityTypePortal
}

func (g *Game) drawPassage(passage world.Passage, field [][]EntityType) {
	for i := 0; i < len(passage.Passage); i++ {
		field[passage.Passage[i].Y][passage.Passage[i].X] = EntityTypePassage
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = EntityTypeDoorOne
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = EntityTypeDoorTwo
}
