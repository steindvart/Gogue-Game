package world

import (
	"fmt"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type DoorPlacer struct{}

func NewDoorPlacer() *DoorPlacer {
	return &DoorPlacer{}
}

func (p *DoorPlacer) PlaceDoors(
	roomOne, roomTwo *Room,
	doorOne, doorTwo primitives.Point2D[int],
) {
	roomOne.Doors = append(roomOne.Doors, doorOne)
	roomTwo.Doors = append(roomTwo.Doors, doorTwo)
}

type ConnectionTreePassageGenerator struct {
	graphBuilder *RoomGraphBuilder
	doorPlacer   *DoorPlacer
}

func NewConnectionTreePassageGenerator() *ConnectionTreePassageGenerator {
	return &ConnectionTreePassageGenerator{
		graphBuilder: NewRoomGraphBuilder(),
		doorPlacer:   NewDoorPlacer(),
	}
}

func (g *ConnectionTreePassageGenerator) GeneratePassages(
	rooms []Room,
	config LevelConfig,
	random utils.Randomizer,
) ([]Passage, error) {
	if len(rooms) < config.RoomsCount {
		return nil, fmt.Errorf("number of rooms is less than expected. expected %d, got %d",
			config.RoomsCount, len(rooms))
	}

	roomIndex := random.Intn(config.RoomsCount)
	treeEdges, err := g.graphBuilder.BuildConnectionTree(roomIndex, config.RoomsCount, random)
	if err != nil {
		return nil, err
	}

	// Добавляем дополнительные рёбра
	extraEdgesCount := random.Intn(config.MaxExtraPassageCount) + 1
	treeEdges = g.graphBuilder.AddRandomEdges(treeEdges, extraEdgesCount, config.RoomsCount, random)

	// Создаём изменяемые копии комнат для добавления дверей
	roomsRef := make([]*Room, len(rooms))
	for i := range rooms {
		roomsRef[i] = &rooms[i]
	}

	// Создаём коридоры и двери
	passages := make([]Passage, 0, len(treeEdges))
	for _, edge := range treeEdges {
		passage, err := g.createPassageForEdge(edge, roomsRef, random)
		if err != nil {
			return nil, err
		}
		passages = append(passages, passage)
	}

	return passages, nil
}

// createPassageForEdge создаёт коридор для данного ребра графа
func (g *ConnectionTreePassageGenerator) createPassageForEdge(
	edge [2]int,
	rooms []*Room,
	random utils.Randomizer,
) (Passage, error) {
	minIndex, maxIndex := sortByOrderAsc(edge[0], edge[1])
	key := [2]int{minIndex, maxIndex}

	// Горизонтальное соединение
	if g.isHorizontalConnection(key) {
		doorOne := getDoorRightWall(*rooms[minIndex], random)
		doorTwo := getDoorLeftWall(*rooms[maxIndex], random)

		passage, err := NewPassageOnX(doorOne, doorTwo, random)
		if err != nil {
			return Passage{}, err
		}

		// Добавляем двери в комнаты
		g.doorPlacer.PlaceDoors(rooms[minIndex], rooms[maxIndex], doorOne, doorTwo)

		return *passage, nil
	}

	// Вертикальное соединение
	if g.isVerticalConnection(key) {
		doorOne := getDoorDownWall(*rooms[minIndex], random)
		doorTwo := getDoorTopWall(*rooms[maxIndex], random)

		passage, err := NewPassageOnY(doorOne, doorTwo, random)
		if err != nil {
			return Passage{}, err
		}

		// Добавляем двери в комнаты
		g.doorPlacer.PlaceDoors(rooms[minIndex], rooms[maxIndex], doorOne, doorTwo)

		return *passage, nil
	}

	return Passage{}, fmt.Errorf("invalid room connection: %v", key)
}

func (g *ConnectionTreePassageGenerator) isHorizontalConnection(key [2]int) bool {
	_, exists := horizontalNeighborRoomsSet[key]
	return exists
}

func (g *ConnectionTreePassageGenerator) isVerticalConnection(key [2]int) bool {
	_, exists := verticalNeighborRoomsSet[key]
	return exists
}

func getDoorLeftWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorYFrom := room.Box.Point.Y + wall
	doorYTo := room.Box.Point.Y + int(room.Box.Size.Height) - wall
	return primitives.Point2D[int]{X: room.Box.Point.X, Y: random.Intn(doorYTo-doorYFrom) + doorYFrom}
}

func getDoorRightWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorYFrom := room.Box.Point.Y + wall
	doorYTo := room.Box.Point.Y + int(room.Box.Size.Height) - wall
	return primitives.Point2D[int]{X: room.Box.Point.X + int(room.Box.Size.Width) - 1, Y: random.Intn(doorYTo-doorYFrom) + doorYFrom}
}

func getDoorTopWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorXFrom := room.Box.Point.X + wall
	doorXTo := room.Box.Point.X + int(room.Box.Size.Width) - wall - 1
	insideRoomWidth := doorXTo - doorXFrom
	var doorX int
	if insideRoomWidth == 0 {
		doorX = doorXFrom
	} else {
		doorX = random.Intn(doorXTo-doorXFrom) + doorXFrom
	}
	return primitives.Point2D[int]{X: doorX, Y: room.Box.Point.Y}
}

func getDoorDownWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorXFrom := room.Box.Point.X + wall
	doorXTo := room.Box.Point.X + int(room.Box.Size.Width) - wall - 1
	insideRoomWidth := doorXTo - doorXFrom
	var doorX int
	if insideRoomWidth == 0 {
		doorX = doorXFrom
	} else {
		doorX = random.Intn(doorXTo-doorXFrom) + doorXFrom
	}
	return primitives.Point2D[int]{X: doorX, Y: room.Box.Point.Y + int(room.Box.Size.Height) - 1}
}
