package world

import (
	"math/rand"
	"testing"
)

func TestRoomGraphBuilder_BuildConnectionTree(t *testing.T) {
	tests := []struct {
		name          string
		startRoom     int
		roomCount     int
		wantTreeEdges [][2]int
	}{
		{
			name:          "Start room 0",
			startRoom:     0,
			roomCount:     9,
			wantTreeEdges: [][2]int{{0, 3}, {0, 1}, {1, 4}, {1, 2}, {2, 5}, {5, 8}, {8, 7}, {7, 6}},
		},
		{
			name:          "Start room 4 (center)",
			startRoom:     4,
			roomCount:     9,
			wantTreeEdges: [][2]int{{4, 5}, {4, 7}, {4, 1}, {4, 3}, {3, 6}, {3, 0}, {1, 2}, {7, 8}},
		},
		{
			name:          "Start room 8",
			startRoom:     8,
			roomCount:     9,
			wantTreeEdges: [][2]int{{8, 7}, {8, 5}, {5, 4}, {5, 2}, {2, 1}, {1, 0}, {0, 3}, {3, 6}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := rand.New(rand.NewSource(testSeed))
			builder := NewRoomGraphBuilder()
			treeEdges, err := builder.BuildConnectionTree(tt.startRoom, tt.roomCount, source)

			if err != nil {
				t.Fatalf("BuildConnectionTree returned unexpected error: %v", err)
			}

			if len(treeEdges) != len(tt.wantTreeEdges) {
				t.Errorf("Expected %d edges in tree, got %d", len(tt.wantTreeEdges), len(treeEdges))
				t.Logf("Real treeEdges for startRoom %d: %v", tt.startRoom, treeEdges)
				return
			}

			for i, wantEdge := range tt.wantTreeEdges {
				if i >= len(treeEdges) {
					t.Errorf("Edge at index %d does not exist in result", i)
					continue
				}
				gotEdge := treeEdges[i]
				if gotEdge != wantEdge {
					t.Errorf("Edge %d: expected %v, got %v", i, wantEdge, gotEdge)
				}
			}
		})
	}
}

func TestRoomGraphBuilder_BuildConnectionTree_Errors(t *testing.T) {
	tests := []struct {
		name      string
		startRoom int
		roomCount int
		wantErr   bool
	}{
		{
			name:      "Start room is negative",
			startRoom: -1,
			roomCount: 9,
			wantErr:   true,
		},
		{
			name:      "Start room exceeds room count",
			startRoom: 12,
			roomCount: 9,
			wantErr:   true,
		},
		{
			name:      "Start room is zero - valid",
			startRoom: 0,
			roomCount: 9,
			wantErr:   false,
		},
		{
			name:      "Start room is last room - valid",
			startRoom: 8,
			roomCount: 9,
			wantErr:   false,
		},
	}

	source := rand.New(rand.NewSource(testSeed))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewRoomGraphBuilder()
			_, err := builder.BuildConnectionTree(tt.startRoom, tt.roomCount, source)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestRoomGraphBuilder_BuildConnectionTree_Connectivity(t *testing.T) {
	// Проверяем что дерево связывает все комнаты
	source := rand.New(rand.NewSource(testSeed))
	builder := NewRoomGraphBuilder()
	roomCount := 9

	for startRoom := 0; startRoom < roomCount; startRoom++ {
		treeEdges, err := builder.BuildConnectionTree(startRoom, roomCount, source)
		if err != nil {
			t.Fatalf("BuildConnectionTree failed for startRoom=%d: %v", startRoom, err)
		}

		// Проверяем что количество рёбер = roomCount - 1 (свойство дерева)
		expectedEdges := roomCount - 1
		if len(treeEdges) != expectedEdges {
			t.Errorf("For startRoom=%d: expected %d edges, got %d", startRoom, expectedEdges, len(treeEdges))
		}

		// Проверяем что все комнаты достижимы
		reachable := make(map[int]bool)
		reachable[startRoom] = true

		// Многократно проходим по рёбрам, добавляя достижимые комнаты
		changed := true
		for changed {
			changed = false
			for _, edge := range treeEdges {
				if reachable[edge[0]] && !reachable[edge[1]] {
					reachable[edge[1]] = true
					changed = true
				}
				if reachable[edge[1]] && !reachable[edge[0]] {
					reachable[edge[0]] = true
					changed = true
				}
			}
		}

		if len(reachable) != roomCount {
			t.Errorf("For startRoom=%d: only %d/%d rooms are reachable", startRoom, len(reachable), roomCount)
		}
	}
}

func TestRoomGraphBuilder_AddRandomEdges(t *testing.T) {
	tests := []struct {
		name            string
		sourceEdges     [][2]int
		extraEdgesCount int
		wantEdgeCount   int
	}{
		{
			name:            "Add zero extra edges",
			sourceEdges:     [][2]int{{0, 1}},
			extraEdgesCount: 0,
			wantEdgeCount:   1 + 1, // sourceEdges + минимум 1
		},
		{
			name:            "Add positive extra edges",
			sourceEdges:     [][2]int{},
			extraEdgesCount: 2,
			wantEdgeCount:   2,
		},
		{
			name:            "Add more extra edges than possible",
			sourceEdges:     [][2]int{},
			extraEdgesCount: 20,
			wantEdgeCount:   2, // максимум 2 для сетки 3x3
		},
		{
			name:            "Add extra edges when some connections already exist",
			sourceEdges:     [][2]int{{0, 1}, {1, 2}},
			extraEdgesCount: 3,
			wantEdgeCount:   2 + 2,
		},
		{
			name:            "Add negative extra edges (should clamp to 1)",
			sourceEdges:     [][2]int{{0, 1}, {1, 2}, {2, 3}},
			extraEdgesCount: -5,
			wantEdgeCount:   3 + 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := rand.New(rand.NewSource(testSeed))
			builder := NewRoomGraphBuilder()
			result := builder.AddRandomEdges(tt.sourceEdges, tt.extraEdgesCount, 9, source)

			if len(result) != tt.wantEdgeCount {
				t.Errorf("Expected %d edges, got %d. Input sourceEdges: %v, extraEdgesCount: %d",
					tt.wantEdgeCount, len(result), tt.sourceEdges, tt.extraEdgesCount)
			}

			// Проверяем что все исходные рёбра сохранились
			resultMap := make(map[[2]int]bool)
			for _, edge := range result {
				minIdx := edge[0]
				maxIdx := edge[1]
				if minIdx > maxIdx {
					minIdx, maxIdx = maxIdx, minIdx
				}
				resultMap[[2]int{minIdx, maxIdx}] = true
			}

			for _, originalEdge := range tt.sourceEdges {
				minIdx := originalEdge[0]
				maxIdx := originalEdge[1]
				if minIdx > maxIdx {
					minIdx, maxIdx = maxIdx, minIdx
				}
				originalKey := [2]int{minIdx, maxIdx}
				if !resultMap[originalKey] {
					t.Errorf("Original edge %v (key %v) is missing from the result", originalEdge, originalKey)
				}
			}
		})
	}
}

func TestRoomGraphBuilder_AddRandomEdges_NoDuplicates(t *testing.T) {
	// Проверяем что AddRandomEdges не создаёт дубликатов
	source := rand.New(rand.NewSource(testSeed))
	builder := NewRoomGraphBuilder()
	sourceEdges := [][2]int{{0, 1}, {1, 2}}
	extraEdgesCount := 5

	result := builder.AddRandomEdges(sourceEdges, extraEdgesCount, 9, source)

	// Проверяем уникальность рёбер
	edgeSet := make(map[[2]int]bool)
	for _, edge := range result {
		minIdx := edge[0]
		maxIdx := edge[1]
		if minIdx > maxIdx {
			minIdx, maxIdx = maxIdx, minIdx
		}
		key := [2]int{minIdx, maxIdx}

		if edgeSet[key] {
			t.Errorf("Duplicate edge found: %v", edge)
		}
		edgeSet[key] = true
	}
}
