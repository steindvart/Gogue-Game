package world

import (
	"errors"
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
	Shape primitives.Box
	Type  RoomType
	Doors []primitives.Point2D[int]

	// occupiedPositions используется только во время генерации для отслеживания занятых позиций
	occupiedPositions map[primitives.Point2D[int]]bool
}

func NewRoom(roomType RoomType, shape primitives.Box) *Room {
	const roomMinWidth = 3
	const roomMinHeight = 3

	if shape.Size.Width < roomMinWidth || shape.Size.Height < roomMinHeight {
		shape.Size.Width = roomMinWidth
		shape.Size.Height = roomMinHeight
	}
	return &Room{
		Shape:             shape,
		Type:              roomType,
		Doors:             []primitives.Point2D[int]{},
		occupiedPositions: make(map[primitives.Point2D[int]]bool),
	}
}

func (r *Room) GetRandomFreePosition(rand utils.Randomizer) (*primitives.Point2D[int], error) {
	// Получаем все возможные точки внутри комнаты без границ
	minX := r.Shape.Point.X + 1
	minY := r.Shape.Point.Y + 1
	width := r.Shape.Size.Width - 2
	height := r.Shape.Size.Height - 2

	totalPossiblePoints := int(width) * int(height)

	for attempt := 0; attempt < totalPossiblePoints; attempt++ {
		pos := primitives.Point2D[int]{
			X: minX + rand.Intn(int(width)),
			Y: minY + rand.Intn(int(height)),
		}

		if !r.occupiedPositions[pos] {
			r.occupiedPositions[pos] = true
			return &pos, nil
		}
	}

	return nil, errors.New("no available positions in room")
}

func (r *Room) MarkOccupied(pos primitives.Point2D[int]) {
	r.occupiedPositions[pos] = true
}

func (r *Room) GetCountFreePosition() int {
	return int(r.Shape.Size.Width*r.Shape.Size.Height) - len(r.occupiedPositions)
}
