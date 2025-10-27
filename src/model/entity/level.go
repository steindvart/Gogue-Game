package entity

import (
	"errors"
	"math"
	"math/rand"
)

const (
	RoomMinWidth     = 3
	RoomMinHeight    = 3
	MinRoomPadding   = 1
	numberXYSections = 3
)

type Level struct {
	Rooms    []Room
	Passages []Passage
	Number   uint
	End      Box
}

var gridNeighborsRooms = [][]int{
	{1, 3},
	{0, 2, 4},
	{1, 5},
	{0, 4, 6},
	{1, 3, 5, 7},
	{2, 4, 8},
	{3, 7},
	{6, 4, 8},
	{5, 7},
}

var horizontalNeighborRoomsSet = map[[2]int]struct{}{
	{0, 1}: {},
	{1, 2}: {},
	{3, 4}: {},
	{4, 5}: {},
	{6, 7}: {},
	{7, 8}: {},
}

var verticalNeighborRoomsSet = map[[2]int]struct{}{
	{0, 3}: {},
	{3, 6}: {},
	{1, 4}: {},
	{4, 7}: {},
	{2, 5}: {},
	{5, 8}: {},
}

func (l *Level) GenerateNineRooms(sizeMap Size2D[uint]) error {
	sectionSize := calculateRoomSectionSize(sizeMap)
	if sectionSize.Width < RoomMinWidth || sectionSize.Height < RoomMinHeight {
		return errors.New("map size is too small: each room section must be at least min room size")
	}

	const roomsCount = 9
	l.Rooms = make([]Room, roomsCount)
	// rand.Perm(9) возвращает массив перемешанных чисел от 0 до 8, чтобы далее не было повторений index для Start и Finish
	indexes := rand.Perm(roomsCount)
	startRoomIndex := indexes[0]
	finishRoomIndex := indexes[1]
	for y := 0; y < numberXYSections; y++ {
		for x := 0; x < numberXYSections; x++ {
			roomIndex := y*numberXYSections + x
			var roomType RoomType
			switch roomIndex {
			case startRoomIndex:
				roomType = RoomTypeStart
			case finishRoomIndex:
				roomType = RoomTypeFinish
			default:
				roomType = RoomTypeOrdinary
			}

			cellXStart := x*int(sectionSize.Width) + MinRoomPadding
			cellYStart := y*int(sectionSize.Height) + MinRoomPadding
			cellXEnd := (x+1)*int(sectionSize.Width) - MinRoomPadding
			cellYEnd := (y+1)*int(sectionSize.Height) - MinRoomPadding

			maxRoomWidth := cellXEnd - cellXStart
			maxRoomHeight := cellYEnd - cellYStart

			if maxRoomWidth < RoomMinWidth || maxRoomHeight < RoomMinHeight {
				return errors.New("internal error: available space in section is smaller than minimum room size")
			}

			width := rand.Intn(maxRoomWidth-RoomMinWidth+1) + RoomMinWidth
			height := rand.Intn(maxRoomHeight-RoomMinHeight+1) + RoomMinHeight

			xCell := cellXStart + rand.Intn(maxRoomWidth-width+1)
			yCell := cellYStart + rand.Intn(maxRoomHeight-height+1)

			if roomType == RoomTypeFinish {
				roomWall := 1
				l.End = Box{
					Point: Point2D[int]{
						X: xCell + roomWall + rand.Intn(width-(roomWall*2)),
						Y: yCell + roomWall + rand.Intn(height-(roomWall*2)),
					},
					Size: Size2D[uint]{Height: 1, Width: 1},
				}
			}

			roomBox := Box{
				Point: Point2D[int]{X: xCell, Y: yCell},
				Size:  Size2D[uint]{Width: uint(width), Height: uint(height)},
			}

			l.Rooms[roomIndex] = *NewRoom(roomType, roomBox)
		}
	}

	return nil
}

func calculateRoomSectionSize(mapSize Size2D[uint]) (sectionSize Size2D[uint]) {
	totalPaddingWidth := uint(MinRoomPadding * 2)
	totalPaddingHeight := uint(MinRoomPadding * 2)

	availableWidth := mapSize.Width - totalPaddingWidth
	availableHeight := mapSize.Height - totalPaddingHeight

	sectionSize.Width = availableWidth / numberXYSections
	sectionSize.Height = availableHeight / numberXYSections

	return sectionSize
}

func (l *Level) GeneratePassages() error {
	roomIndex := rand.Intn(9)
	treeEdges, err := generateSpanningTree(roomIndex)
	if err != nil {
		return err
	}

	extraEdgesCount := rand.Intn(4)
	treeEdges = addRandomEdges(treeEdges, extraEdgesCount)

	for _, connectedRooms := range treeEdges {
		minIndex, maxIndex := sortByOrderAsc(connectedRooms[0], connectedRooms[1])
		key := [2]int{minIndex, maxIndex}
		if _, exists := horizontalNeighborRoomsSet[key]; exists {
			doorOne := getDoorRightWall(l.Rooms[minIndex])
			doorTwo := getDoorLeftWall(l.Rooms[maxIndex])
			l.Passages = append(l.Passages, *NewPassageOnX(doorOne, doorTwo))
		}
		if _, exists := verticalNeighborRoomsSet[key]; exists {
			doorOne := getDoorDownWall(l.Rooms[minIndex])
			doorTwo := getDoorTopWall(l.Rooms[maxIndex])
			l.Passages = append(l.Passages, *NewPassageOnY(doorOne, doorTwo))
		}
	}

	return nil
}

func generateSpanningTree(startRoom int) ([][2]int, error) {
	if startRoom < 0 || startRoom > 9 {
		return nil, errors.New("start room must be between 0 and 9")
	}

	edges := make([][2]int, 0)
	stack := []int{startRoom}
	visit := make(map[int]bool)
	visit[startRoom] = true

	for len(stack) > 0 {
		currentRoom := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		allNeighborsRooms := gridNeighborsRooms[currentRoom]
		unvisitNeighborsRooms := []int{}
		for _, neighborRoom := range allNeighborsRooms {
			if !visit[neighborRoom] {
				unvisitNeighborsRooms = append(unvisitNeighborsRooms, neighborRoom)
			}
		}

		// rand.Shuffle перемешивает значения в существующем слайсе unvisitNeighborsRooms
		rand.Shuffle(len(unvisitNeighborsRooms), func(i, j int) {
			unvisitNeighborsRooms[i], unvisitNeighborsRooms[j] = unvisitNeighborsRooms[j], unvisitNeighborsRooms[i]
		})

		for _, nextRoom := range unvisitNeighborsRooms {
			edges = append(edges, [2]int{currentRoom, nextRoom})
			visit[nextRoom] = true

			stack = append(stack, nextRoom)
		}
	}

	return edges, nil
}

func addRandomEdges(sourceEdges [][2]int, extraEdgesCount int) [][2]int {
	existingConnections := make(map[[2]int]struct{})
	for _, connectedRooms := range sourceEdges {
		minIndex, maxIndex := sortByOrderAsc(connectedRooms[0], connectedRooms[1])
		key := [2]int{minIndex, maxIndex}
		existingConnections[key] = struct{}{}
	}

	unionSet := make(map[[2]int]struct{})

	for k := range horizontalNeighborRoomsSet {
		unionSet[k] = struct{}{}
	}

	for k := range verticalNeighborRoomsSet {
		unionSet[k] = struct{}{}
	}

	resultSet := make(map[[2]int]struct{})

	for k := range unionSet {
		if _, existsInAnother := existingConnections[k]; !existsInAnother {
			resultSet[k] = struct{}{}
		}
	}

	potentialEdges := make([][2]int, len(resultSet))
	for k := range resultSet {
		potentialEdges = append(potentialEdges, k)
	}

	rand.Shuffle(len(potentialEdges), func(i, j int) {
		potentialEdges[i], potentialEdges[j] = potentialEdges[j], potentialEdges[i]
	})

	resultEdges := make([][2]int, len(sourceEdges))
	for i := range sourceEdges {
		resultEdges[i] = sourceEdges[i]
	}

	for i := 0; i < extraEdgesCount; i++ {
		element := potentialEdges[i]
		resultEdges = append(resultEdges, [2]int{element[0], element[1]})
	}

	return resultEdges
}

func getDoorLeftWall(room Room) Point2D[int] {
	const wall = 1
	doorYFrom := room.Shape.Point.Y + wall
	doorYTo := room.Shape.Point.Y + int(room.Shape.Size.Height) - wall
	return Point2D[int]{X: room.Shape.Point.X, Y: rand.Intn(doorYTo-doorYFrom) + doorYFrom}
}

func getDoorRightWall(room Room) Point2D[int] {
	const wall = 1
	doorYFrom := room.Shape.Point.Y + wall
	doorYTo := room.Shape.Point.Y + int(room.Shape.Size.Height) - wall
	return Point2D[int]{X: room.Shape.Point.X + int(room.Shape.Size.Width), Y: rand.Intn(doorYTo-doorYFrom) + doorYFrom}
}

func getDoorTopWall(room Room) Point2D[int] {
	const wall = 1
	doorXFrom := room.Shape.Point.X + wall
	doorXTo := room.Shape.Point.X + int(room.Shape.Size.Width) - wall
	return Point2D[int]{X: rand.Intn(doorXTo-doorXFrom) + doorXFrom, Y: room.Shape.Point.Y}
}

func getDoorDownWall(room Room) Point2D[int] {
	const wall = 1
	doorXFrom := room.Shape.Point.X + wall
	doorXTo := room.Shape.Point.X + int(room.Shape.Size.Width)
	return Point2D[int]{X: rand.Intn(doorXTo-doorXFrom) + doorXFrom, Y: room.Shape.Point.Y + int(room.Shape.Size.Height)}
}

func sortByOrderAsc(first, second int) (int, int) {
	minIndex := int(math.Min(float64(first), float64(second)))
	maxIndex := int(math.Max(float64(first), float64(second)))
	return minIndex, maxIndex
}
