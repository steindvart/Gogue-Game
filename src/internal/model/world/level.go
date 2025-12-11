package world

import (
	"errors"
	"fmt"
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"math"
	"sort"
)

const (
	RoomMinWidth         = 3
	RoomMinHeight        = 3
	MapMaxWidth          = 200
	MapMaxHeight         = 150
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
	Player       *entities.Player
	Number       uint
	FinishPortal primitives.Box
	random       utils.Randomizer
	mapSize      primitives.Size2D[uint]
}

func NewLevel(random utils.Randomizer, mapSize primitives.Size2D[uint]) *Level {
	return &Level{
		random:  random,
		mapSize: mapSize,
	}
}

func (l *Level) Generate() error {
	err := l.generateMap()
	if err != nil {
		return err
	}

	err = l.generatePlayer()
	if err != nil {
		return err
	}

	return nil
}

func (l *Level) GenerateWithExistingPlayer(player *entities.Player) error {
	err := l.generateMap()
	if err != nil {
		return err
	}

	l.Player = player

	startPlayerPos, err := l.GenerateStartPlayerPosition()
	if err != nil {
		return err
	}

	l.Player.SetPosition(*startPlayerPos)

	return nil
}

func (l *Level) generateMap() error {
	err := l.generateNineRooms()
	if err != nil {
		return err
	}

	err = l.generatePassages()
	if err != nil {
		return err
	}

	// @todo - "С каждым новым уровнем снижается количество полезных предметов (и повышается количество сокровищ, которые выпадают с побежденных противников)"
	// конкретные значения будут задаваться не тут, но это готовые "ручки", которые дёрнутся, когда будет готова логика уровня
	err = l.addItemsAtRooms(7, 5, 3, 1)
	if err != nil {
		return err
	}

	return nil
}

func (l *Level) generatePlayer() error {
	startPlayerPos, err := l.GenerateStartPlayerPosition()
	if err != nil {
		return err
	}

	l.Player = entities.NewPlayer(&primitives.Box{
		Point: *startPlayerPos,
		Size:  primitives.Size2D[uint]{Height: 1, Width: 1},
	})

	return nil
}

func (l *Level) generateNineRooms() error {
	sectionSize, err := calculateRoomSectionSize(l.mapSize)
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
			roomType := RoomTypeOrdinary
			switch roomIndex {
			case startRoomIndex:
				roomType = RoomTypeStart
			case finishRoomIndex:
				roomType = RoomTypeFinish
			}

			cellXStart := x*int(sectionSize.Width) + MinRoomPadding
			cellYStart := y*int(sectionSize.Height) + MinRoomPadding
			cellXEnd := (x + 1) * int(sectionSize.Width)
			cellYEnd := (y + 1) * int(sectionSize.Height)

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
				portalPos, errPos := l.Rooms[roomIndex].GetRandomFreePosition(l.random)
				if errPos != nil {
					return errors.New("no free positions in finish room to place level portal")
				}
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
	if len(l.Rooms) < roomsCount {
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

func (l *Level) addItemsAtRooms(countFood, countElixir, countScroll, countWeapon uint) error {
	if l.Rooms == nil {
		return errors.New("no rooms on level")
	}

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
	allScrollTypes := []items.ScrollType{
		items.ScrollTypeStrength,
		items.ScrollTypeAgility,
		items.ScrollTypeUltimate,
		items.ScrollTypeMaxHealth,
		items.ScrollTypeMystery,
		// items.ScrollTypeCustom,
	}
	allWeaponTypes := []items.WeaponType{
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
			return nil
		}
		l.Rooms[roomInd].createFood(l.random, *pos, allFoodTypes)
	}
	for j := 0; j < int(countElixir); j++ {
		roomInd, pos, err := l.getFreePosition()
		if err != nil {
			return nil
		}
		l.Rooms[roomInd].createElixir(l.random, *pos, allElixirTypes)
	}
	for j := 0; j < int(countScroll); j++ {
		roomInd, pos, err := l.getFreePosition()
		if err != nil {
			return nil
		}
		l.Rooms[roomInd].createScroll(l.random, *pos, allScrollTypes)
	}
	for j := 0; j < int(countWeapon); j++ {
		roomInd, pos, err := l.getFreePosition()
		if err != nil {
			return nil
		}
		l.Rooms[roomInd].createWeapon(l.random, *pos, allWeaponTypes)
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

		pos, err := l.Rooms[ind].GetRandomFreePosition(l.random)
		if err != nil {
			continue
		}

		return ind, pos, nil
	}

	return 0, nil, errors.New("no available positions in Rooms")
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
			pos, err := room.GetRandomFreePosition(l.random)
			if err != nil {
				return nil, fmt.Errorf("in starting room not free position")
			}
			return pos, nil
		}
	}

	return nil, fmt.Errorf("starting room was not found")
}

func calculateRoomSectionSize(sizeMap primitives.Size2D[uint]) (primitives.Size2D[uint], error) {
	if sizeMap.Width > MapMaxWidth || sizeMap.Height > MapMaxHeight {
		return primitives.Size2D[uint]{}, errors.New("game map size is too big")
	}

	availableWidth := sizeMap.Width
	availableHeight := sizeMap.Height

	sectionSize := primitives.Size2D[uint]{
		Width:  availableWidth / numberXYSections,
		Height: availableHeight / numberXYSections,
	}

	const areaForPassage = 1
	if sectionSize.Width < RoomMinWidth+areaForPassage || sectionSize.Height < RoomMinHeight+areaForPassage {
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
	return primitives.Point2D[int]{X: room.Shape.Point.X + int(room.Shape.Size.Width) - 1, Y: random.Intn(doorYTo-doorYFrom) + doorYFrom}
}

func getDoorTopWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorXFrom := room.Shape.Point.X + wall
	doorXTo := room.Shape.Point.X + int(room.Shape.Size.Width) - wall - 1
	insideRoomWidth := doorXTo - doorXFrom
	var doorX int
	if insideRoomWidth == 0 {
		doorX = doorXFrom
	} else {
		doorX = random.Intn(doorXTo-doorXFrom) + doorXFrom
	}
	return primitives.Point2D[int]{X: doorX, Y: room.Shape.Point.Y}
}

func getDoorDownWall(room Room, random utils.Randomizer) primitives.Point2D[int] {
	const wall = 1
	doorXFrom := room.Shape.Point.X + wall
	doorXTo := room.Shape.Point.X + int(room.Shape.Size.Width) - wall - 1
	insideRoomWidth := doorXTo - doorXFrom
	var doorX int
	if insideRoomWidth == 0 {
		doorX = doorXFrom
	} else {
		doorX = random.Intn(doorXTo-doorXFrom) + doorXFrom
	}
	return primitives.Point2D[int]{X: doorX, Y: room.Shape.Point.Y + int(room.Shape.Size.Height) - 1}
}

func (l *Level) MovePlayer(delta primitives.Point2D[int]) {
	oldPlayerPos := l.Player.GetPosition()
	l.Player.Move(delta)

	if l.checkCollision(l.Player.GetPosition()) {
		l.Player.SetPosition(oldPlayerPos)
	}

}

func (l *Level) checkCollision(pos primitives.Point2D[int]) bool {
	if l.checkCollisionWithMapBorders(pos) {
		return true
	}
	if l.checkCollisionWithRoomsWall(pos) {
		return true
	}
	if l.checkCollisionWithEnemy(pos) {
		return true
	}

	if !isInSomeRoom(pos, l.Rooms) && !l.checkCollisionWithPassages(pos) {
		return true
	}

	return false
}

func (l *Level) checkCollisionWithMapBorders(pos primitives.Point2D[int]) bool {
	if pos.X < 0 || pos.Y < 0 {
		return true
	}
	if pos.X >= int(l.mapSize.Width) || pos.Y >= int(l.mapSize.Height) {
		return true
	}

	return false
}

func (l *Level) checkCollisionWithRoomsWall(pos primitives.Point2D[int]) bool {
	for _, room := range l.Rooms {
		if isInRoom(pos, room) && checkCollisionWithRoomWall(pos, room) {
			return true
		}
	}

	return false
}

func isInSomeRoom(pos primitives.Point2D[int], rooms []Room) bool {
	for _, room := range rooms {
		if isInRoom(pos, room) {
			return true
		}
	}

	return false
}

func isInRoom(pos primitives.Point2D[int], room Room) bool {
	// Проверка того, что персонаж находится внутри области комнаты
	leftEndX := room.Shape.Point.X + 1
	rightEndX := leftEndX + int(room.Shape.Size.Width) - 3
	topEndY := room.Shape.Point.Y + 1
	downEndY := topEndY + int(room.Shape.Size.Height) - 3

	return (pos.X >= leftEndX && pos.X <= rightEndX) &&
		(pos.Y >= topEndY && pos.Y <= downEndY)
}

func checkCollisionWithRoomWall(pos primitives.Point2D[int], room Room) bool {
	// Двери являются частью комнаты и её стен, но через них можно ходить
	if checkCollisionWithDoors(pos, room.Doors) {
		return false
	}

	leftEndX := room.Shape.Point.X
	rightEndX := leftEndX + int(room.Shape.Size.Width)
	topEndY := room.Shape.Point.Y
	downEndY := topEndY + int(room.Shape.Size.Height)

	if (pos.X == leftEndX || pos.X == rightEndX) || (pos.Y == topEndY || pos.Y == downEndY) {
		return true
	}

	return false
}

func checkCollisionWithDoors(pos primitives.Point2D[int], doors []primitives.Point2D[int]) bool {
	for _, door := range doors {
		if pos == door {
			return true
		}
	}

	return false
}

func (l *Level) checkCollisionWithPassages(newPos primitives.Point2D[int]) bool {
	for _, passage := range l.Passages {
		if isInPassage(newPos, passage) {
			return true
		}
	}

	return false
}

func isInPassage(pos primitives.Point2D[int], passage Passage) bool {
	// Двери являются как частью комнаты, так и частью прохода
	if pos == passage.DoorOne || pos == passage.DoorTwo {
		return true
	}

	for _, wayPoint := range passage.Way {
		if pos == wayPoint {
			return true
		}
	}

	return false
}

func (l *Level) checkCollisionWithEnemy(pos primitives.Point2D[int]) bool {
	for _, room := range l.Rooms {
		for _, enemy := range room.Enemies {
			if pos == enemy.GetPosition() {
				return true
			}
		}
	}

	return false
}

func (l *Level) MakeCurrentMap(w, h int) [][]common.GameEntityType {
	field := make([][]common.GameEntityType, h)
	for y := range field {
		field[y] = make([]common.GameEntityType, w)
	}

	// Добавлено в качестве примера, потом надо будет убрать
	l.Rooms[0].Enemies = []entities.Enemy{
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeZombie),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 6},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeVampire),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 7},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeGhost),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 8},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeOgre),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 9},
		//		},
		//	},
		//},
		//{
		//	Type: entities.EnemyType(entities.EnemyTypeSnakeMage),
		//	Character: &entities.Character{
		//		Shape: &primitives.Box{
		//			Point: primitives.Point2D[int]{X: 6, Y: 10},
		//		},
		//	},
		//},
	}

	for _, room := range l.Rooms {
		l.putRoom(room, l.FinishPortal, field)
		l.putItems(room, field)
	}

	for _, passages := range l.Passages {
		l.putPassage(passages, field)
	}

	// for _, room := range l.Rooms {
	// 	g.putEnemies(room, field)
	// }

	px := l.Player.Character.Shape.Point.X
	py := l.Player.Character.Shape.Point.Y
	if py >= 0 && py < h && px >= 0 && px < w {
		field[py][px] = common.EntityTypePlayer
	}

	return field
}

//func (g *Game) putEnemies(room Room, field [][]common.GameEntityType) {
//	var et common.GameEntityType
//
//	for _, e := range room.Enemies {
//		switch e.Type {
//		case entities.EnemyTypeZombie:
//			et = common.EntityTypeZombie
//		case entities.EnemyTypeVampire:
//			et = common.EntityTypeVampire
//		case entities.EnemyTypeGhost:
//			et = common.EntityTypeGhost
//		case entities.EnemyTypeOgre:
//			et = common.EntityTypeOgre
//		case entities.EnemyTypeSnakeMage:
//			et = common.EntityTypeSnakeMage
//		}
//
//		ex := e.Character.Shape.Point.X
//		ey := e.Character.Shape.Point.Y
//
//		h := len(field)
//		w := len(field[0])
//		if ey >= 0 && ey < h && ex >= 0 && ex < w {
//			field[ey][ex] = et
//		}
//	}
//}

func (l *Level) putRoom(room Room, finishPortal primitives.Box, field [][]common.GameEntityType) {
	width := int(room.Shape.Size.Width)
	height := int(room.Shape.Size.Height)

	startX := room.Shape.Point.X
	endX := startX + width - 1
	startY := room.Shape.Point.Y
	endY := startY + height - 1

	for col := startX; col <= endX; col++ {
		field[startY][col] = common.WorldTypeWall
		field[endY][col] = common.WorldTypeWall
	}

	for row := startY; row <= endY; row++ {
		field[row][startX] = common.WorldTypeWall
		field[row][endX] = common.WorldTypeWall
	}

	field[finishPortal.Point.Y][finishPortal.Point.X] = common.WorldTypePortal
}

func (l *Level) putPassage(passage Passage, field [][]common.GameEntityType) {
	for i := 0; i < len(passage.Way); i++ {
		field[passage.Way[i].Y][passage.Way[i].X] = common.WorldTypePassage
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = common.WorldTypeDoor
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = common.WorldTypeDoor
}

// @todo - оставить ли это всё кучей?
func (l *Level) putItems(room Room, field [][]common.GameEntityType) {
	var fd common.GameEntityType
	for _, food := range room.Foods {
		switch food.Type {
		case items.FoodTypePotatoes:
			fd = common.FoodTypePotatoes
		case items.FoodTypeBread:
			fd = common.FoodTypeBread
		case items.FoodTypeMeat:
			fd = common.FoodTypeMeat
		case items.FoodTypeMistery:
			fd = common.FoodTypeMistery
		case items.FoodTypeBeer:
			fd = common.FoodTypeBeer
		}
		field[food.Item.Shape.Point.Y][food.Item.Shape.Point.X] = fd
	}

	for _, elixir := range room.Elixirs {
		switch elixir.Type {
		case items.ElixirTypeStrength:
			fd = common.ElixirTypeStrength
		case items.ElixirTypeAgility:
			fd = common.ElixirTypeAgility
		case items.ElixirTypeDwarfism:
			fd = common.ElixirTypeDwarfism
		case items.ElixirTypeGiantism:
			fd = common.ElixirTypeGiantism
		case items.ElixirTypeMystery:
			fd = common.ElixirTypeMystery
		}
		field[elixir.Item.Shape.Point.Y][elixir.Item.Shape.Point.X] = fd
	}

	for _, scroll := range room.Scrolls {
		switch scroll.Type {
		case items.ScrollTypeStrength:
			fd = common.ScrollTypeStrength
		case items.ScrollTypeAgility:
			fd = common.ScrollTypeAgility
		case items.ScrollTypeUltimate:
			fd = common.ScrollTypeUltimate
		case items.ScrollTypeMaxHealth:
			fd = common.ScrollTypeMaxHealth
		case items.ScrollTypeMystery:
			fd = common.ScrollTypeMystery
		}
		field[scroll.Item.Shape.Point.Y][scroll.Item.Shape.Point.X] = fd
	}

	for _, weapon := range room.Weapons {
		fd = common.Weapon
		field[weapon.Item.Shape.Point.Y][weapon.Item.Shape.Point.X] = fd
	}
}
