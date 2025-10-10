package state

import (
	"fmt"
	"gogue/model/entity"
	"gogue/model/signal"
	"gogue/view/action"
	view "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

type Game struct {
	player entity.Player
	view   *view.Game
}

func NewGame(parent *gc.Window) *Game {
	viewObj, err := view.NewGame(parent)
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
		player: player,
		view:   viewObj,
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
	field := g.generateField()
	g.view.Render(field)
}

// Генерирует двумерное поле, где 0 — пусто, 1 — персонаж
func (g *Game) generateField() [][]int {
	// Размер поля (например, 20x20)
	width, height := 20, 20
	field := make([][]int, height)
	for y := range field {
		field[y] = make([]int, width)
	}
	// Координаты персонажа
	px := g.player.Character.Shape.Point.X
	py := g.player.Character.Shape.Point.Y
	// Проверка границ
	if py >= 0 && py < height && px >= 0 && px < width {
		field[py][px] = 1
	}
	return field
}
