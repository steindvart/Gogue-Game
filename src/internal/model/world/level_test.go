package world

import (
	"errors"
	"gogue/internal/model/primitives"
	"math/rand"
	"strings"
	"testing"
)

const randomSeedTest = 42

func TestLevel_GenerateNineRooms(t *testing.T) {
	tests := []struct {
		name      string
		sizeMap   primitives.Size2D[uint]
		wantErr   bool
		errorText string
	}{
		{
			name:    "Valid map size 30x90",
			sizeMap: primitives.Size2D[uint]{Height: 30, Width: 90},
			wantErr: false,
		},
		{
			name:      "Map too small",
			sizeMap:   primitives.Size2D[uint]{Height: 6, Width: 6},
			wantErr:   true,
			errorText: "game map size is too small",
		},
		{
			name:      "Map too big",
			sizeMap:   primitives.Size2D[uint]{Height: 300, Width: 300},
			wantErr:   true,
			errorText: "game map size is too big",
		},
	}

	source := rand.New(rand.NewSource(randomSeedTest))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := NewLevel(source)

			err := level.generateNineRooms(tt.sizeMap)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tt.errorText)
				}
				if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLevel_GenerateNineRooms_Deterministic(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	testSize := primitives.Size2D[uint]{Height: 30, Width: 90}

	level := NewLevel(source)
	err := level.generateNineRooms(testSize)
	if err != nil {
		t.Fatalf("generateNineRooms failed: %v", err)
	}

	wantRooms := []Room{
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 4, Y: 1},
				Size:  primitives.Size2D[uint]{Width: 21, Height: 7},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 35, Y: 3},
				Size:  primitives.Size2D[uint]{Width: 6, Height: 5},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 67, Y: 3},
				Size:  primitives.Size2D[uint]{Width: 4, Height: 3},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 12, Y: 10},
				Size:  primitives.Size2D[uint]{Width: 12, Height: 7},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 39, Y: 11},
				Size:  primitives.Size2D[uint]{Width: 14, Height: 5},
			},
		},
		{
			Type: RoomTypeFinish,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 61, Y: 11},
				Size:  primitives.Size2D[uint]{Width: 7, Height: 5},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 2, Y: 19},
				Size:  primitives.Size2D[uint]{Width: 26, Height: 7},
			},
		},
		{
			Type: RoomTypeStart,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 37, Y: 23},
				Size:  primitives.Size2D[uint]{Width: 10, Height: 3},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 66, Y: 21},
				Size:  primitives.Size2D[uint]{Width: 13, Height: 4},
			},
		},
	}

	wantFinishPortal := primitives.Box{
		Point: primitives.Point2D[int]{X: 62, Y: 14},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
	if level.FinishPortal != wantFinishPortal {
		t.Errorf("Expected FinishPortal %+v, got %+v", wantFinishPortal, level.FinishPortal)
	}

	if len(level.Rooms) != len(wantRooms) {
		t.Errorf("Expected %d rooms, got %d", len(wantRooms), len(level.Rooms))
		t.Logf("Generated Rooms: %+v", level.Rooms)
		return
	}

	for i, wantRoom := range wantRooms {
		if i >= len(level.Rooms) {
			t.Errorf("Room at index %d does not exist in generated level", i)
			continue
		}
		gotRoom := level.Rooms[i]

		if gotRoom.Type != wantRoom.Type {
			t.Errorf("Room %d: expected Type %v, got %v", i, wantRoom.Type, gotRoom.Type)
		}
		if gotRoom.Shape != wantRoom.Shape {
			t.Errorf("Room %d: expected Shape %+v, got %+v", i, wantRoom.Shape, gotRoom.Shape)
		}
	}
}

func TestLevel_GeneratePassages_Deterministic(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	tests := []struct {
		name         string
		wantPassages []Passage
	}{
		{
			name: "Generate Passages with deterministic rooms",
			wantPassages: []Passage{
				{
					DoorOne: primitives.Point2D[int]{X: 25, Y: 3},
					DoorTwo: primitives.Point2D[int]{X: 35, Y: 4},
					Passage: []primitives.Point2D[int]{
						{X: 26, Y: 3}, {X: 27, Y: 3}, {X: 28, Y: 3}, {X: 29, Y: 3}, {X: 30, Y: 3},
						{X: 30, Y: 4}, {X: 31, Y: 4}, {X: 32, Y: 4}, {X: 33, Y: 4}, {X: 34, Y: 4},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 9, Y: 8},
					DoorTwo: primitives.Point2D[int]{X: 21, Y: 10},
					Passage: []primitives.Point2D[int]{
						{X: 9, Y: 9}, {X: 10, Y: 9}, {X: 11, Y: 9}, {X: 12, Y: 9},
						{X: 13, Y: 9}, {X: 14, Y: 9}, {X: 15, Y: 9}, {X: 16, Y: 9},
						{X: 17, Y: 9}, {X: 18, Y: 9}, {X: 19, Y: 9}, {X: 20, Y: 9}, {X: 21, Y: 9},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 24, Y: 12},
					DoorTwo: primitives.Point2D[int]{X: 39, Y: 12},
					Passage: []primitives.Point2D[int]{
						{X: 25, Y: 12}, {X: 26, Y: 12}, {X: 27, Y: 12}, {X: 28, Y: 12}, {X: 29, Y: 12},
						{X: 30, Y: 12}, {X: 31, Y: 12}, {X: 32, Y: 12}, {X: 33, Y: 12}, {X: 34, Y: 12},
						{X: 35, Y: 12}, {X: 36, Y: 12}, {X: 37, Y: 12}, {X: 38, Y: 12},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 13, Y: 17},
					DoorTwo: primitives.Point2D[int]{X: 11, Y: 19},
					Passage: []primitives.Point2D[int]{
						{X: 13, Y: 18}, {X: 12, Y: 18}, {X: 11, Y: 18},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 28, Y: 20},
					DoorTwo: primitives.Point2D[int]{X: 37, Y: 24},
					Passage: []primitives.Point2D[int]{
						{X: 29, Y: 20}, {X: 30, Y: 20}, {X: 31, Y: 20}, {X: 31, Y: 21},
						{X: 31, Y: 22}, {X: 31, Y: 23}, {X: 31, Y: 24}, {X: 32, Y: 24},
						{X: 33, Y: 24}, {X: 34, Y: 24}, {X: 35, Y: 24}, {X: 36, Y: 24},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 47, Y: 24},
					DoorTwo: primitives.Point2D[int]{X: 66, Y: 23},
					Passage: []primitives.Point2D[int]{
						{X: 48, Y: 24}, {X: 49, Y: 24}, {X: 50, Y: 24}, {X: 51, Y: 24}, {X: 52, Y: 24},
						{X: 53, Y: 24}, {X: 54, Y: 24}, {X: 55, Y: 24}, {X: 56, Y: 24}, {X: 57, Y: 24},
						{X: 57, Y: 23}, {X: 58, Y: 23}, {X: 59, Y: 23}, {X: 60, Y: 23}, {X: 61, Y: 23},
						{X: 62, Y: 23}, {X: 63, Y: 23}, {X: 64, Y: 23}, {X: 65, Y: 23},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 66, Y: 16},
					DoorTwo: primitives.Point2D[int]{X: 69, Y: 21},
					Passage: []primitives.Point2D[int]{
						{X: 66, Y: 17}, {X: 67, Y: 17}, {X: 68, Y: 17},
						{X: 69, Y: 17}, {X: 69, Y: 18}, {X: 69, Y: 19}, {X: 69, Y: 20},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 68, Y: 6},
					DoorTwo: primitives.Point2D[int]{X: 63, Y: 11},
					Passage: []primitives.Point2D[int]{
						{X: 68, Y: 7}, {X: 68, Y: 8}, {X: 68, Y: 9}, {X: 67, Y: 9},
						{X: 66, Y: 9}, {X: 65, Y: 9}, {X: 64, Y: 9}, {X: 63, Y: 9}, {X: 63, Y: 10},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 39, Y: 8},
					DoorTwo: primitives.Point2D[int]{X: 41, Y: 11},
					Passage: []primitives.Point2D[int]{
						{X: 39, Y: 9}, {X: 40, Y: 9}, {X: 41, Y: 9}, {X: 41, Y: 10},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := NewLevel(source)
			err := level.generateNineRooms(primitives.Size2D[uint]{Height: 30, Width: 90})
			if err != nil {
				t.Fatalf("generateNineRooms returned an error: %v", err)
			}
			err = level.generatePassages()
			if err != nil {
				t.Fatalf("generatePassages returned an error: %v", err)
			}

			if len(level.Passages) != len(tt.wantPassages) {
				t.Fatalf("Length mismatch: expected %d passages, got %d passages.",
					len(tt.wantPassages), len(level.Passages))
			}

			for i := range tt.wantPassages {
				gotPassage := level.Passages[i]
				wantPassage := tt.wantPassages[i]

				if gotPassage.DoorOne != wantPassage.DoorOne {
					t.Errorf("Passage at index %d field 'DoorOne' mismatches: expected %v, got %v",
						i, wantPassage.DoorOne, gotPassage.DoorOne)
				}
				if gotPassage.DoorTwo != wantPassage.DoorTwo {
					t.Errorf("Passage at index %d field 'DoorOne' mismatches: expected %v, got %v",
						i, wantPassage.DoorTwo, gotPassage.DoorTwo)
				}
				if len(gotPassage.Passage) != len(wantPassage.Passage) {
					t.Errorf("Passage at index %d Passage slice length mismatch: expected %d, got %d", i, len(wantPassage.Passage), len(gotPassage.Passage))
				} else {
					for j, wantPoint := range wantPassage.Passage {
						gotPoint := gotPassage.Passage[j]
						if gotPoint != wantPoint {
							t.Errorf("Passage at index %d, point in Passage at index %d mismatch: expected %+v, got %+v", i, j, wantPoint, gotPoint)
						}
					}
				}
			}
		})
	}
}

func TestLevel_CalculateRoomSectionSize(t *testing.T) {
	tests := []struct {
		name    string
		sizeMap primitives.Size2D[uint]
		want    primitives.Size2D[uint]
		wantErr error
	}{
		{
			name:    "Valid size 33x33",
			sizeMap: primitives.Size2D[uint]{Width: 33, Height: 33},
			want:    primitives.Size2D[uint]{Width: 10, Height: 10},
			wantErr: nil,
		},
		{
			name:    "Valid size 60x45",
			sizeMap: primitives.Size2D[uint]{Width: 60, Height: 45},
			want:    primitives.Size2D[uint]{Width: 19, Height: 14},
			wantErr: nil,
		},
		{
			name:    "Size too small for rooms",
			sizeMap: primitives.Size2D[uint]{Width: 10, Height: 10},
			want:    primitives.Size2D[uint]{},
			wantErr: errors.New("game map size is too small"),
		},
		{
			name:    "Size too big - width",
			sizeMap: primitives.Size2D[uint]{Width: 201, Height: 100},
			want:    primitives.Size2D[uint]{},
			wantErr: errors.New("game map size is too big"),
		},
		{
			name:    "Size too big - height",
			sizeMap: primitives.Size2D[uint]{Width: 100, Height: 151},
			want:    primitives.Size2D[uint]{},
			wantErr: errors.New("game map size is too big"),
		},
		{
			name:    "Size too big - both",
			sizeMap: primitives.Size2D[uint]{Width: 300, Height: 200},
			want:    primitives.Size2D[uint]{},
			wantErr: errors.New("game map size is too big"),
		},
		{
			name:    "Minimum valid size 15x15",
			sizeMap: primitives.Size2D[uint]{Width: 15, Height: 15},
			want:    primitives.Size2D[uint]{Width: 4, Height: 4},
			wantErr: nil,
		},
		{
			name:    "Size just below minimum - 10x10",
			sizeMap: primitives.Size2D[uint]{Width: 10, Height: 10},
			want:    primitives.Size2D[uint]{},
			wantErr: errors.New("game map size is too small"),
		},
		{
			name:    "Size with remainder - 35x34",
			sizeMap: primitives.Size2D[uint]{Width: 35, Height: 34},
			want:    primitives.Size2D[uint]{Width: 11, Height: 10},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := calculateRoomSectionSize(tt.sizeMap)

			if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error to contain: %v, got: %v", tt.wantErr, err)
				}
			} else if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error to contain: %v, got: <nil>", tt.wantErr)
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected: <nil>, got error: %v", err)
			}

			if tt.wantErr == nil {
				if res != tt.want {
					t.Errorf("Expected size (H, W): %v, got size (H, W): %v", tt.want, res)
				}
			}
		})
	}
}

func TestLevel_GenerateSpanningTree_Errors(t *testing.T) {
	tests := []struct {
		name      string
		startRoom int
		wantErr   error
	}{
		{
			name:      "Start room is negative",
			startRoom: -1,
			wantErr:   errors.New("start room must be between 0 and 9"),
		},
		{
			name:      "Start room is too big",
			startRoom: 12,
			wantErr:   errors.New("start room must be between 0 and 9"),
		},
		{
			name:      "Start room is zero",
			startRoom: 0,
			wantErr:   nil,
		},
		{
			name:      "Start room is eight",
			startRoom: 8,
			wantErr:   nil,
		},
	}

	source := rand.New(rand.NewSource(randomSeedTest))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generateSpanningTree(tt.startRoom, source)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("Expected error containing %v, but got <nil>", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Expected Success for startRoom=%d, but got error: %v", tt.startRoom, err)
				}
			}
		})
	}
}

func TestLevel_GenerateSpanningTree_Deterministic(t *testing.T) {
	tests := []struct {
		name          string
		startRoom     int
		wantTreeEdges [][2]int
	}{
		{
			name:          "Start room 0",
			startRoom:     0,
			wantTreeEdges: [][2]int{{0, 3}, {0, 1}, {1, 4}, {1, 2}, {2, 5}, {5, 8}, {8, 7}, {7, 6}},
		},
		{
			name:          "Start room 4",
			startRoom:     4,
			wantTreeEdges: [][2]int{{4, 5}, {4, 7}, {4, 1}, {4, 3}, {3, 6}, {3, 0}, {1, 2}, {7, 8}},
		},
		{
			name:          "Start room 8",
			startRoom:     8,
			wantTreeEdges: [][2]int{{8, 7}, {8, 5}, {5, 4}, {5, 2}, {2, 1}, {1, 0}, {0, 3}, {3, 6}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := rand.New(rand.NewSource(randomSeedTest))
			treeEdges, err := generateSpanningTree(tt.startRoom, source)

			if err != nil {
				t.Fatalf("generateSpanningTree returned unexpected error: %v", err)
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

func TestLevel_AddRandomEdges(t *testing.T) {
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
			wantEdgeCount:   1 + 1,
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
			wantEdgeCount:   2,
		},
		{
			name:            "Add extra edges when some connections already exist",
			sourceEdges:     [][2]int{{0, 1}, {1, 2}},
			extraEdgesCount: 3,
			wantEdgeCount:   2 + 2,
		},
		{
			name:            "Add negative extra edges",
			sourceEdges:     [][2]int{{0, 1}, {1, 2}, {2, 3}},
			extraEdgesCount: -5,
			wantEdgeCount:   3 + 1,
		},
	}

	source := rand.New(rand.NewSource(randomSeedTest))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addRandomEdges(tt.sourceEdges, tt.extraEdgesCount, source)

			if len(result) != tt.wantEdgeCount {
				t.Errorf("Expected %d edges, got %d. Input sourceEdges: %v, extraEdgesCount: %d", tt.wantEdgeCount, len(result), tt.sourceEdges, tt.extraEdgesCount)
			}

			resultMap := make(map[[2]int]bool)
			for _, edge := range result {
				minIdx, maxIdx := sortByOrderAsc(edge[0], edge[1])
				resultMap[[2]int{minIdx, maxIdx}] = true
			}
			for _, originalEdge := range tt.sourceEdges {
				minIdx, maxIdx := sortByOrderAsc(originalEdge[0], originalEdge[1])
				originalKey := [2]int{minIdx, maxIdx}
				if !resultMap[originalKey] {
					t.Errorf("Original edge %v (key %v) is missing from the result", originalEdge, originalKey)
				}
			}
		})
	}
}

func TestLevel_AddRandomEdges_Deterministic(t *testing.T) {
	tests := []struct {
		name            string
		sourceEdges     [][2]int
		extraEdgesCount int
		wantResult      [][2]int
	}{
		{
			name:            "Add 1 extra to empty",
			sourceEdges:     [][2]int{{0, 3}, {0, 1}},
			extraEdgesCount: 1,
			wantResult:      [][2]int{{0, 3}, {0, 1}, {6, 7}},
		},
		{
			name:            "Add 2 extra to simple tree",
			sourceEdges:     [][2]int{{0, 3}, {0, 1}},
			extraEdgesCount: 2,
			wantResult:      [][2]int{{0, 3}, {0, 1}, {6, 7}, {4, 7}},
		},
		{
			name:            "Add -1 extra (clamped to 1)",
			sourceEdges:     [][2]int{{0, 3}, {0, 1}},
			extraEdgesCount: -1,
			wantResult:      [][2]int{{0, 3}, {0, 1}, {6, 7}},
		},
		{
			name:            "Add 25 extra (clamped to 1)",
			sourceEdges:     [][2]int{{0, 3}, {0, 1}},
			extraEdgesCount: 25,
			wantResult:      [][2]int{{0, 3}, {0, 1}, {6, 7}, {4, 7}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := rand.New(rand.NewSource(randomSeedTest))
			result := addRandomEdges(tt.sourceEdges, tt.extraEdgesCount, source)

			if len(result) != len(tt.wantResult) {
				t.Fatalf("Length mismatch: expected %d edges, got %d edges.\nReal result for sourceEdges=%v, extraCount=%d: %v",
					len(tt.wantResult), len(result), tt.sourceEdges, tt.extraEdgesCount, result)
			}

			for i := range tt.wantResult {
				if result[i] != tt.wantResult[i] {
					t.Errorf("Edge at index %d mismatches: expected %v, got %v", i, tt.wantResult[i], result[i])
					t.Logf("Full result: %v", result)
					t.Logf("Full wantResult: %v", tt.wantResult)
				}
			}
		})
	}
}

func TestLevel_GetStartPositionForPlayer(t *testing.T) {
	t.Run("Rooms exist", func(t *testing.T) {
		wantPlayerStartPoint := primitives.Point2D[int]{X: 40, Y: 24}

		source := rand.New(rand.NewSource(randomSeedTest))
		level := NewLevel(source)
		err := level.GenerateLevel(primitives.Size2D[uint]{Height: 30, Width: 90})
		if err != nil {
			t.Fatalf("GenerateLevel returned unexpected error: %v", err)
		}
		playerStartPoint, err := level.GetStartPositionForPlayer()
		if err != nil {
			t.Fatalf("GetStartPositionForPlayer returned unexpected error: %v", err)
		}

		if wantPlayerStartPoint != *playerStartPoint {
			t.Errorf("Position mismatch: expected %+v, got %+v", wantPlayerStartPoint, playerStartPoint)
		}
	})
	t.Run("Rooms don't exist", func(t *testing.T) {
		wantErrorText := "should be 9 rooms"
		source := rand.New(rand.NewSource(randomSeedTest))
		level := NewLevel(source)
		_, err := level.GetStartPositionForPlayer()
		if err == nil {
			t.Fatalf("GetStartPositionForPlayer returned nil error, expected %q", wantErrorText)
		}
		if !strings.Contains(err.Error(), wantErrorText) {
			t.Errorf("Expected error to contain %q, got %q", wantErrorText, err.Error())
		}
	})
}
