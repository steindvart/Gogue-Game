package world

import (
	"fmt"
	"gogue/internal/utils"
	"math"
	"sort"
)

// gridNeighborsRooms содержит список соседей для каждой комнаты в сетке 3x3
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

// horizontalNeighborRoomsSet содержит пары комнат, соединяемых горизонтально
var horizontalNeighborRoomsSet = map[[2]int]struct{}{
	{0, 1}: {},
	{1, 2}: {},
	{3, 4}: {},
	{4, 5}: {},
	{6, 7}: {},
	{7, 8}: {},
}

// verticalNeighborRoomsSet содержит пары комнат, соединяемых вертикально
var verticalNeighborRoomsSet = map[[2]int]struct{}{
	{0, 3}: {},
	{3, 6}: {},
	{1, 4}: {},
	{4, 7}: {},
	{2, 5}: {},
	{5, 8}: {},
}

// RoomGraphBuilder строит граф связности комнат
type RoomGraphBuilder struct{}

func NewRoomGraphBuilder() *RoomGraphBuilder {
	return &RoomGraphBuilder{}
}

func (b *RoomGraphBuilder) BuildConnectionTree(startRoom, roomCount int, random utils.Randomizer) ([][2]int, error) {
	if startRoom < 0 || startRoom >= roomCount {
		return nil, fmt.Errorf("start room must be between 0 and %d", roomCount-1)
	}

	edges := make([][2]int, 0)
	stack := []int{startRoom}
	visited := make(map[int]bool)
	visited[startRoom] = true

	for len(stack) > 0 {
		currentRoom := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		neighbors := b.getUnvisitedNeighbors(currentRoom, visited)
		random.Shuffle(len(neighbors), func(i, j int) {
			neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
		})

		for _, nextRoom := range neighbors {
			edges = append(edges, [2]int{currentRoom, nextRoom})
			visited[nextRoom] = true
			stack = append(stack, nextRoom)
		}
	}

	return edges, nil
}

func (b *RoomGraphBuilder) getUnvisitedNeighbors(room int, visited map[int]bool) []int {
	allNeighbors := gridNeighborsRooms[room]
	unvisited := make([]int, 0)

	for _, neighbor := range allNeighbors {
		if !visited[neighbor] {
			unvisited = append(unvisited, neighbor)
		}
	}

	return unvisited
}

func (b *RoomGraphBuilder) AddRandomEdges(
	sourceEdges [][2]int,
	extraCount int,
	roomCount int,
	random utils.Randomizer,
) [][2]int {
	const minExtraCount = 1
	const maxExtraCount = 2

	if extraCount < minExtraCount {
		extraCount = minExtraCount
	}
	if extraCount > maxExtraCount {
		extraCount = maxExtraCount
	}

	existingConnections := b.buildConnectionSet(sourceEdges)
	potentialEdges := b.findPotentialEdges(existingConnections)

	if len(potentialEdges) == 0 {
		return sourceEdges
	}

	// Перемешиваем дважды для большей случайности
	random.Shuffle(len(potentialEdges), func(i, j int) {
		potentialEdges[i], potentialEdges[j] = potentialEdges[j], potentialEdges[i]
	})
	random.Shuffle(len(potentialEdges), func(i, j int) {
		potentialEdges[i], potentialEdges[j] = potentialEdges[j], potentialEdges[i]
	})

	resultEdges := make([][2]int, len(sourceEdges))
	copy(resultEdges, sourceEdges)

	for i := 0; i < extraCount && i < len(potentialEdges); i++ {
		resultEdges = append(resultEdges, potentialEdges[i])
	}

	return resultEdges
}

func (b *RoomGraphBuilder) buildConnectionSet(edges [][2]int) map[[2]int]struct{} {
	connections := make(map[[2]int]struct{})
	for _, edge := range edges {
		minIdx, maxIdx := sortByOrderAsc(edge[0], edge[1])
		connections[[2]int{minIdx, maxIdx}] = struct{}{}
	}
	return connections
}

func (b *RoomGraphBuilder) findPotentialEdges(existingConnections map[[2]int]struct{}) [][2]int {
	allPossible := b.getAllPossibleConnections()
	potential := make([][2]int, 0)

	for conn := range allPossible {
		if _, exists := existingConnections[conn]; !exists {
			potential = append(potential, conn)
		}
	}

	// Сортируем для детерминизма в тестах
	sort.Slice(potential, func(i, j int) bool {
		if potential[i][0] != potential[j][0] {
			return potential[i][0] < potential[j][0]
		}
		return potential[i][1] < potential[j][1]
	})

	return potential
}

func (b *RoomGraphBuilder) getAllPossibleConnections() map[[2]int]struct{} {
	allConnections := make(map[[2]int]struct{})

	for k := range horizontalNeighborRoomsSet {
		allConnections[k] = struct{}{}
	}
	for k := range verticalNeighborRoomsSet {
		allConnections[k] = struct{}{}
	}

	return allConnections
}

// sortByOrderAsc возвращает два индекса в возрастающем порядке
func sortByOrderAsc(first, second int) (int, int) {
	return int(math.Min(float64(first), float64(second))),
		int(math.Max(float64(first), float64(second)))
}
