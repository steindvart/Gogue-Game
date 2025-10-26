package entity

import (
	"errors"
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

var xPassages = [][2]int{
	{0, 1},
	{1, 2},
	{3, 4},
	{4, 5},
	{6, 7},
	{7, 8},
}

var yPassages = [][2]int{
	{0, 3},
	{3, 6},
	{1, 4},
	{4, 7},
	{2, 5},
	{5, 8},
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

func generatePassagesTree(startRoom int) (map[int][]int, error) {
	if startRoom < 0 || startRoom > 9 {
		return nil, errors.New("start room must be between 0 and 9")
	}
	graph := make(map[int][]int)
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

		// rand.Shuffle перемешивает значения в слайсе unvisitNeighborsRooms
		rand.Shuffle(len(unvisitNeighborsRooms), func(i, j int) {
			unvisitNeighborsRooms[i], unvisitNeighborsRooms[j] = unvisitNeighborsRooms[j], unvisitNeighborsRooms[i]
		})

		for _, nextRoom := range unvisitNeighborsRooms {
			graph[currentRoom] = append(graph[currentRoom], nextRoom)
			graph[nextRoom] = append(graph[nextRoom], currentRoom)
			visit[nextRoom] = true

			stack = append(stack, nextRoom)
		}
	}

	return graph, nil
}

func (l *Level) GeneratePassages() error {
	tree, err := generatePassagesTree(0)
	if err != nil {
		return err
	}

	// [2]int - две комнаты, между которыми тоннель. bool - наличие тоннеля.
	laidPassages := make(map[[2]int]bool)
	for roomOne, connectedRooms := range tree {
		for _, roomTwo := range connectedRooms {
			var key [2]int
			if roomOne < roomTwo {
				key = [2]int{roomOne, roomTwo}
			} else {
				key = [2]int{roomTwo, roomOne}
			}
			if !laidPassages[key] {
				laidPassages[key] = true
				for _, passageKey := range xPassages {
					if key == passageKey {
						indexOne := passageKey[0]
						indexTwo := passageKey[1]
						doorOne := getDoorRightWall(l.Rooms[indexOne])
						doorTwo := getDoorLeftWall(l.Rooms[indexTwo])
						l.Passages = append(l.Passages, *NewPassageOnX(doorOne, doorTwo))
					}
				}
				for _, passageKey := range yPassages {
					if key == passageKey {
						indexOne := passageKey[0]
						indexTwo := passageKey[1]
						doorOne := getDoorDownWall(l.Rooms[indexOne])
						doorTwo := getDoorTopWall(l.Rooms[indexTwo])
						l.Passages = append(l.Passages, *NewPassageOnY(doorOne, doorTwo))
					}
				}
			}
		}
	}

	return nil
}
