package world

import (
	"errors"
	"fmt"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"math"
	"sort"
)

const (
	RoomMinWidth         = 3
	RoomMinHeight        = 3
	RoomMaxWidth         = 200
	RoomMaxHeight        = 150
	MinRoomPadding       = 1
	MaxExtraPassageCount = 2
	numberXYSections     = 3
	roomsCount           = 9
)

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

type Level struct {
	Rooms        []Room
	Passages     []Passage
	Number       uint
	FinishPortal primitives.Box
	random       utils.Randomizer
}

func NewLevel(random utils.Randomizer) *Level {
	return &Level{
		random: random,
	}
}

func (l *Level) GenerateLevel(sizeMap primitives.Size2D[uint]) error {
	err := l.generateNineRooms(sizeMap)
	if err != nil {
		return err
	}

	err = l.generatePassages()
	if err != nil {
		return err
	}

	// @todo - "С каждым новым уровнем снижается количество полезных предметов (и повышается количество сокровищ, которые выпадают с побежденных противников)"
	l.addItemsAtRooms(2, 2, 2, 2)

	return nil
}

func (l *Level) generateNineRooms(sizeMap primitives.Size2D[uint]) error {
	sectionSize, err := calculateRoomSectionSize(sizeMap)
	if err != nil {
		return err
	}

	l.Rooms = make([]Room, roomsCount)
	// rand.Perm(9) возвращает массив перемешанных чисел от 0 до 8
	indexes := l.random.Perm(roomsCount)
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
				return errors.New("available space in section is smaller than minimum room size")
			}

			width := l.random.Intn(maxRoomWidth-RoomMinWidth+1) + RoomMinWidth
			height := l.random.Intn(maxRoomHeight-RoomMinHeight+1) + RoomMinHeight

			xCell := cellXStart + l.random.Intn(maxRoomWidth-width+1)
			yCell := cellYStart + l.random.Intn(maxRoomHeight-height+1)

			roomBox := primitives.Box{
				Point: primitives.Point2D[int]{X: xCell, Y: yCell},
				Size:  primitives.Size2D[uint]{Width: uint(width), Height: uint(height)},
			}

			l.Rooms[roomIndex] = *NewRoom(roomType, roomBox)

			if roomType == RoomTypeFinish {
				portalPos := l.Rooms[roomIndex].GetRandomFreePosition(l.random)
				l.FinishPortal = primitives.Box{
					Point: *portalPos,
					Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
				}
			}
		}
	}

	return nil
}

func (l *Level) generatePassages() error {
	if len(l.Rooms) < 9 {
		return fmt.Errorf("number of rooms is less than expected. expected %d, got %d", roomsCount, len(l.Rooms))
	}
	roomIndex := l.random.Intn(roomsCount)
	treeEdges, err := generateSpanningTree(roomIndex, l.random)
	if err != nil {
		return err
	}

	extraEdgesCount := l.random.Intn(MaxExtraPassageCount) + 1
	treeEdges = addRandomEdges(treeEdges, extraEdgesCount, l.random)

	for _, connectedRooms := range treeEdges {
		minIndex, maxIndex := sortByOrderAsc(connectedRooms[0], connectedRooms[1])
		key := [2]int{minIndex, maxIndex}

		if _, exists := horizontalNeighborRoomsSet[key]; exists {
			doorOne := getDoorRightWall(l.Rooms[minIndex], l.random)
			doorTwo := getDoorLeftWall(l.Rooms[maxIndex], l.random)
			passage, err := NewPassageOnX(doorOne, doorTwo, l.random)
			if err != nil {
				return err
			}
			l.Passages = append(l.Passages, *passage)
			keyAsUint := [2]uint{uint(key[0]), uint(key[1])}
			err = l.addDoorsAtRoom(keyAsUint, doorOne, doorTwo)
			if err != nil {
				return err
			}
		}
		if _, exists := verticalNeighborRoomsSet[key]; exists {
			doorOne := getDoorDownWall(l.Rooms[minIndex], l.random)
			doorTwo := getDoorTopWall(l.Rooms[maxIndex], l.random)
			passage, err := NewPassageOnY(doorOne, doorTwo, l.random)
			if err != nil {
				return err
			}
			l.Passages = append(l.Passages, *passage)
			keyAsUint := [2]uint{uint(key[0]), uint(key[1])}
			err = l.addDoorsAtRoom(keyAsUint, doorOne, doorTwo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *Level) getFreePosition() (int, *primitives.Point2D[int], error) {
	indexes := l.random.Perm(roomsCount)

	for indCount := 0; indCount < len(l.Rooms); indCount++ {
		ind := indexes[indCount]
		if l.Rooms[ind].Type == RoomTypeStart {
			continue
		}

		if l.Rooms[ind].GetCountFreePosition() <= 0 {
			continue
		}

		pos := l.Rooms[ind].GetRandomFreePosition(l.random)
		if pos == nil {
			continue
		}

		return ind, pos, nil
	}

	return 0, nil, errors.New("no available positions in Rooms")
}

func (l *Level) addItemsAtRooms(countFood, countElixir, countScroll, countWeapon uint) {
	allFoodTypes := []items.FoodType{
		items.FoodTypePotatoes,
		items.FoodTypeBread,
		items.FoodTypeMeat,
		items.FoodTypeMistery,
		items.FoodTypeBeer,
	}
	allElixirTypes := []items.ElixirType{
		items.ElixirTypeStrength,
		items.ElixirTypeAgility,
		items.ElixirTypeDwarfism,
		items.ElixirTypeGiantism,
		items.ElixirTypeMystery,
		// items.ElixirTypeCustom,
	}
	var allScrollTypes = []items.ScrollType{
		items.ScrollTypeStrength,
		items.ScrollTypeAgility,
		items.ScrollTypeUltimate,
		items.ScrollTypeMaxHealth,
		items.ScrollTypeMystery,
		// items.ScrollTypeCustom,
	}
	var allWeaponTypes = []items.WeaponType{
		items.WeaponTypeDagger,
		items.WeaponTypeSpear,
		items.WeaponTypeSword,
		items.WeaponTypeAxe,
		items.WeaponTypeMaul,
		items.WeaponTypeMystery,
		// items.WeaponTypeCustom,
	}

	for j := 0; j < int(countFood); j++ {
		roomInd, pos, err := l.getFreePosition()
		if err != nil {
			return
		}
		l.createFood(roomInd, *pos, allFoodTypes)
	}
	for j := 0; j < int(countElixir); j++ {
		roomInd, pos, err := l.getFreePosition()
		if err != nil {
			return
		}
		l.createElixir(roomInd, *pos, allElixirTypes)
	}
	for j := 0; j < int(countScroll); j++ {
		roomInd, pos, err := l.getFreePosition()
		if err != nil {
			return
		}
		l.createScroll(roomInd, *pos, allScrollTypes)
	}
	for j := 0; j < int(countWeapon); j++ {
		roomInd, pos, err := l.getFreePosition()
		if err != nil {
			return
		}
		l.createWeapon(roomInd, *pos, allWeaponTypes)
	}
}

func getRandomElement[T any](random utils.Randomizer, slice []T) T {
	idx := random.Intn(len(slice))
	return slice[idx]
}

func (l *Level) createFood(roomInd int, pos primitives.Point2D[int], allFoodType []items.FoodType) {
	foodType := getRandomElement(l.random, allFoodType)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newFood := items.NewFoodBuiltin(l.random, itemBox, foodType)
	l.Rooms[roomInd].Foods = append(l.Rooms[roomInd].Foods, *newFood)
}

func (l *Level) createElixir(roomInd int, pos primitives.Point2D[int], allElixirType []items.ElixirType) bool {
	elixirType := getRandomElement(l.random, allElixirType)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newElixir := items.NewElixirBuiltin(l.random, itemBox, elixirType)
	l.Rooms[roomInd].Elixirs = append(l.Rooms[roomInd].Elixirs, *newElixir)
	return true
}

func (l *Level) createScroll(roomInd int, pos primitives.Point2D[int], allScrollType []items.ScrollType) bool {
	scrollType := getRandomElement(l.random, allScrollType)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newScroll := items.NewScrollBuiltin(l.random, itemBox, scrollType)
	l.Rooms[roomInd].Scrolls = append(l.Rooms[roomInd].Scrolls, *newScroll)
	return true
}

func (l *Level) createWeapon(roomInd int, pos primitives.Point2D[int], allWeaponTypes []items.WeaponType) bool {
	weaponType := getRandomElement(l.random, allWeaponTypes)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newWeapon := items.NewWeaponBuiltin(l.random, itemBox, weaponType)
	l.Rooms[roomInd].Weapons = append(l.Rooms[roomInd].Weapons, *newWeapon)
	return true
}

func (l *Level) addDoorsAtRoom(twoRoomsIndexes [2]uint, doorOne, doorTwo primitives.Point2D[int]) error {
	for _, roomIndex := range twoRoomsIndexes {
		if roomIndex > roomsCount-1 {
			return fmt.Errorf("rooms indexes should be in range from 0 to %d", roomsCount-1)
		}
	}
	roomOneIndex := twoRoomsIndexes[0]
	l.Rooms[roomOneIndex].Doors = append(l.Rooms[roomOneIndex].Doors, doorOne)
	roomTwoIndex := twoRoomsIndexes[1]
	l.Rooms[roomTwoIndex].Doors = append(l.Rooms[roomTwoIndex].Doors, doorTwo)

	return nil
}

func (l *Level) GenerateStartPlayerPosition() (*primitives.Point2D[int], error) {
	if len(l.Rooms) != roomsCount {
		return nil, fmt.Errorf("should be %d rooms", roomsCount)
	}

	for _, room := range l.Rooms {
		if room.Type == RoomTypeStart {
			return room.GetRandomFreePosition(l.random), nil
		}
	}

	return nil, fmt.Errorf("starting room was not found")
}

func calculateRoomSectionSize(sizeMap primitives.Size2D[uint]) (primitives.Size2D[uint], error) {
	if sizeMap.Width > RoomMaxWidth || sizeMap.Height > RoomMaxHeight {
		return primitives.Size2D[uint]{}, errors.New("game map size is too big")
	}

	totalPaddingWidth := uint(MinRoomPadding * 2)
	totalPaddingHeight := uint(MinRoomPadding * 2)

	availableWidth := sizeMap.Width - totalPaddingWidth
	availableHeight := sizeMap.Height - totalPaddingHeight

	sectionSize := primitives.Size2D[uint]{
		Width:  availableWidth / numberXYSections,
		Height: availableHeight / numberXYSections,
	}

	if sectionSize.Width < RoomMinWidth || sectionSize.Height < RoomMinHeight {
		return primitives.Size2D[uint]{}, errors.New("game map size is too small")
	}

	return sectionSize, nil
}

func generateSpanningTree(startRoom int, random utils.Randomizer) ([][2]int, error) {
	if startRoom < 0 || startRoom > roomsCount {
		return nil, fmt.Errorf("start room must be between 0 and %d", roomsCount)
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
		random.Shuffle(len(unvisitNeighborsRooms), func(i, j int) {
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

func addRandomEdges(sourceEdges [][2]int, extraEdgesCount int, random utils.Randomizer) [][2]int {
	const minExtraPassageCount = 1
	if extraEdgesCount < minExtraPassageCount {
		extraEdgesCount = minExtraPassageCount
	}
	if extraEdgesCount > MaxExtraPassageCount {
		extraEdgesCount = MaxExtraPassageCount
	}

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

	potentialEdges := make([][2]int, 0)
	for k := range resultSet {
		potentialEdges = append(potentialEdges, k)
	}

	sort.Slice(potentialEdges, func(i, j int) bool {
		if potentialEdges[i][0] != potentialEdges[j][0] {
			return potentialEdges[i][0] < potentialEdges[j][0]
		}
		return potentialEdges[i][1] < potentialEdges[j][1]
	})

	random.Shuffle(len(potentialEdges), func(i, j int) {
		potentialEdges[i], potentialEdges[j] = potentialEdges[j], potentialEdges[i]
	})

	random.Shuffle(len(potentialEdges), func(i, j int) {
		potentialEdges[i], potentialEdges[j] = potentialEdges[j], potentialEdges[i]
	})

	resultEdges := make([][2]int, len(sourceEdges))
	copy(resultEdges, sourceEdges)

	for i := 0; i < extraEdgesCount; i++ {
		element := potentialEdges[i]
		resultEdges = append(resultEdges, [2]int{element[0], element[1]})
	}

	return resultEdges
}

func sortByOrderAsc(first, second int) (int, int) {
	minIndex := int(math.Min(float64(first), float64(second)))
	maxIndex := int(math.Max(float64(first), float64(second)))
	return minIndex, maxIndex
}

func getDoorLeftWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorYFrom := room.Shape.Point.Y + wall
	doorYTo := room.Shape.Point.Y + int(room.Shape.Size.Height) - wall
	return primitives.Point2D[int]{X: room.Shape.Point.X, Y: random.Intn(doorYTo-doorYFrom) + doorYFrom}
}

func getDoorRightWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorYFrom := room.Shape.Point.Y + wall
	doorYTo := room.Shape.Point.Y + int(room.Shape.Size.Height) - wall
	return primitives.Point2D[int]{X: room.Shape.Point.X + int(room.Shape.Size.Width), Y: random.Intn(doorYTo-doorYFrom) + doorYFrom}
}

func getDoorTopWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorXFrom := room.Shape.Point.X + wall
	doorXTo := room.Shape.Point.X + int(room.Shape.Size.Width) - wall
	return primitives.Point2D[int]{X: random.Intn(doorXTo-doorXFrom) + doorXFrom, Y: room.Shape.Point.Y}
}

func getDoorDownWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorXFrom := room.Shape.Point.X + wall
	doorXTo := room.Shape.Point.X + int(room.Shape.Size.Width)
	return primitives.Point2D[int]{X: random.Intn(doorXTo-doorXFrom) + doorXFrom, Y: room.Shape.Point.Y + int(room.Shape.Size.Height)}
}
