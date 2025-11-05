package entity

import (
	"errors"
	"strings"
	"testing"
)

func TestLevel_GenerateNineRooms(t *testing.T) {
	tests := []struct {
		name      string
		sizeMap   Size2D[uint]
		wantErr   bool
		errorText string
	}{
		{
			name:    "Valid map size 30x90",
			sizeMap: Size2D[uint]{Height: 30, Width: 90},
			wantErr: false,
		},
		{
			name:      "Map too small",
			sizeMap:   Size2D[uint]{Height: 6, Width: 6},
			wantErr:   true,
			errorText: "game map size is too small",
		},
		{
			name:      "Map too big",
			sizeMap:   Size2D[uint]{Height: 300, Width: 300},
			wantErr:   true,
			errorText: "game map size is too big",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := &Level{}

			err := level.GenerateNineRooms(tt.sizeMap)

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

func TestLevel_CalculateRoomSectionSize(t *testing.T) {
	type testCase struct {
		name    string
		sizeMap Size2D[uint]
		want    Size2D[uint]
		wantErr error
	}

	tests := []testCase{
		{
			name:    "Valid size 33x33",
			sizeMap: Size2D[uint]{Width: 33, Height: 33},
			want:    Size2D[uint]{Width: 10, Height: 10},
			wantErr: nil,
		},
		{
			name:    "Valid size 60x45",
			sizeMap: Size2D[uint]{Width: 60, Height: 45},
			want:    Size2D[uint]{Width: 19, Height: 14},
			wantErr: nil,
		},
		{
			name:    "Size too small for rooms",
			sizeMap: Size2D[uint]{Width: 10, Height: 10},
			want:    Size2D[uint]{},
			wantErr: errors.New("game map size is too small"),
		},
		{
			name:    "Size too big - width",
			sizeMap: Size2D[uint]{Width: 201, Height: 100},
			want:    Size2D[uint]{},
			wantErr: errors.New("game map size is too big"),
		},
		{
			name:    "Size too big - height",
			sizeMap: Size2D[uint]{Width: 100, Height: 151},
			want:    Size2D[uint]{},
			wantErr: errors.New("game map size is too big"),
		},
		{
			name:    "Size too big - both",
			sizeMap: Size2D[uint]{Width: 300, Height: 200},
			want:    Size2D[uint]{},
			wantErr: errors.New("game map size is too big"),
		},
		{
			name:    "Minimum valid size 15x15",
			sizeMap: Size2D[uint]{Width: 15, Height: 15},
			want:    Size2D[uint]{Width: 4, Height: 4},
			wantErr: nil,
		},
		{
			name:    "Size just below minimum - 10x10",
			sizeMap: Size2D[uint]{Width: 10, Height: 10},
			want:    Size2D[uint]{},
			wantErr: errors.New("game map size is too small"),
		},
		{
			name:    "Size with remainder - 35x34",
			sizeMap: Size2D[uint]{Width: 35, Height: 34},
			want:    Size2D[uint]{Width: 11, Height: 10},
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
	type testCase struct {
		name      string
		startRoom int
		wantErr   error
	}

	tests := []testCase{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generateSpanningTree(tt.startRoom)

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

func TestLevel_AddRandomEdges(t *testing.T) {
	type testCase struct {
		name            string
		sourceEdges     [][2]int
		extraEdgesCount int
		wantEdgeCount   int
	}
	tests := []testCase{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addRandomEdges(tt.sourceEdges, tt.extraEdgesCount)

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
