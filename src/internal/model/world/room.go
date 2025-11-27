package world

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type RoomType uint

const (
	RoomTypeStart RoomType = iota
	RoomTypeOrdinary
	RoomTypeFinish
)

type Room struct {
	Shape   primitives.Box
	Type    RoomType
	Doors   []primitives.Point2D[int]
	Foods   []items.Food
	Elixirs []items.Elixir
	Scrolls []items.Scroll
	Weapons []items.Weapon
	Enemies []entities.Enemy

	OccupiedPositions map[primitives.Point2D[int]]bool
}

func NewRoom(roomType RoomType, shape primitives.Box) *Room {
	return &Room{
		Shape:             shape,
		Type:              roomType,
		Foods:             []items.Food{},
		Elixirs:           []items.Elixir{},
		Scrolls:           []items.Scroll{},
		Weapons:           []items.Weapon{},
		Enemies:           []entities.Enemy{},
		OccupiedPositions: make(map[primitives.Point2D[int]]bool),
	}
}

func (r *Room) GetRandomFreePosition(rand utils.RandomSource) *primitives.Point2D[int] {
	// Получаем все возможные точки внутри комнаты без границ
	minX := r.Shape.Point.X + 1
	maxX := r.Shape.Point.X + int(r.Shape.Size.Width) - 1
	minY := r.Shape.Point.Y + 1
	maxY := r.Shape.Point.Y + int(r.Shape.Size.Height) - 1

	width := maxX - minX
	height := maxY - minY

	if width <= 0 || height <= 0 {
		return nil
	}

	totalPossiblePoints := width * height

	for attempt := 0; attempt < totalPossiblePoints; attempt++ {
		p := primitives.Point2D[int]{
			X: minX + rand.Intn(width),
			Y: minY + rand.Intn(height),
		}

		if !r.OccupiedPositions[p] {
			r.OccupiedPositions[p] = true
			return &p
		}
	}

	return nil
}
