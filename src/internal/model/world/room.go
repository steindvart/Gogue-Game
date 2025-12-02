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

func (r *Room) GetRandomFreePosition(rand utils.Randomizer) *primitives.Point2D[int] {
	// Получаем все возможные точки внутри комнаты без границ
	minX := r.Shape.Point.X + 1
	minY := r.Shape.Point.Y + 1
	width := r.Shape.Size.Width - 1
	height := r.Shape.Size.Height - 1

	if width <= 0 || height <= 0 {
		return nil
	}

	totalPossiblePoints := int(width) * int(height)

	for attempt := 0; attempt < totalPossiblePoints; attempt++ {
		pos := primitives.Point2D[int]{
			X: minX + rand.Intn(int(width)),
			Y: minY + rand.Intn(int(height)),
		}

		if !r.OccupiedPositions[pos] {
			r.OccupiedPositions[pos] = true
			return &pos
		}
	}

	return nil
}

func (r *Room) GetCountFreePosition() int {
	return int(r.Shape.Size.Width*r.Shape.Size.Height) - len(r.OccupiedPositions)
}
