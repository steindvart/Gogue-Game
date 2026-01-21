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
	*primitives.Box
	Doors []primitives.Point2D[int]

	// Type используется только во время генерации для указания комнаты старта для Player и конечной комнаты для Portal
	Type RoomType
	// occupiedPositions используется только во время генерации для отслеживания занятых позиций
	occupiedPositions map[primitives.Point2D[int]]bool
}

func NewRoom(roomType RoomType, box primitives.Box) (*Room, error) {
	const roomMinWidth = 3
	const roomMinHeight = 3

	if box.Size.Width < roomMinWidth || box.Size.Height < roomMinHeight {
		return nil, errors.New("room size Width and Height should be more 3")
	}
	return &Room{
		Box:               &box,
		Type:              roomType,
		Doors:             []primitives.Point2D[int]{},
		occupiedPositions: make(map[primitives.Point2D[int]]bool),
	}, nil
}

func (r *Room) GetRandomFreePosition(rand utils.Randomizer) (*primitives.Point2D[int], error) {
	// Получаем все возможные точки внутри комнаты без границ
	minX := r.Box.Point.X + 1
	minY := r.Box.Point.Y + 1
	width := r.Box.Size.Width - 2
	height := r.Box.Size.Height - 2

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
	return int(r.Box.Size.Width*r.Box.Size.Height) - len(r.occupiedPositions)
}
