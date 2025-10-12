package state

import (
	"fmt"
	"gogue/model/entity"
	"gogue/model/signal"
	"gogue/view/action"
	view "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

const (
	FIELD_HEIGHT = 40 // размеры с учётом границ
	FIELD_WIDTH  = 100
)

type GameRenderer interface {
	Render(field [][]int) error
}

type Game struct {
	player      entity.Player
	fieldWidth  int
	fieldHeight int
	renderer    GameRenderer
}

func NewGame(parent *gc.Window) *Game {
	renderer, err := view.NewGame(parent, FIELD_WIDTH, FIELD_HEIGHT)
	if err != nil {
		fmt.Println("Error creating main menu view:", err)
		return nil
	}

	player := entity.Player{
		Character: entity.Character{
			Shape: entity.Box{
				Point: entity.Point2D[int]{X: 10, Y: 10},
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

	return &Game{
		player:      player,
		fieldWidth:  FIELD_WIDTH - 2,  // Исключая границы
		fieldHeight: FIELD_HEIGHT - 2, // Исключая границы
		renderer:    renderer,
	}
}

func (g *Game) Input(a action.Type) signal.Type {
	switch a {
	case action.MoveUp:
		g.player.Character.Shape.Move(entity.Point2D[int]{X: 0, Y: -1})
	case action.MoveDown:
		g.player.Character.Shape.Move(entity.Point2D[int]{X: 0, Y: 1})
	case action.MoveLeft:
		g.player.Character.Shape.Move(entity.Point2D[int]{X: -1, Y: 0})
	case action.MoveRight:
		g.player.Character.Shape.Move(entity.Point2D[int]{X: 1, Y: 0})
	case action.Exit:
		return signal.Stop
	}

	return signal.NoSignal
}

func (g *Game) Update() signal.Type {
	return signal.NoSignal
}

func (g *Game) Render() {
	field := g.makeField()

	err := g.renderer.Render(field)
	if err != nil {
		fmt.Println("Error rendering game state:", err)
		panic(err)
	}
}

// Генерирует двумерное поле, где 0 — пусто, 1 — персонаж
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
