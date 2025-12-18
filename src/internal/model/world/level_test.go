package world

import (
	"gogue/internal/model/primitives"
	"math/rand"
	"strings"
	"testing"
)

const randomSeedTest = 42

func TestLevel_Generate(t *testing.T) {
	tests := []struct {
		name      string
		mapSize   primitives.Size2D[uint]
		wantErr   bool
		errorText string
	}{
		{
			name:    "Valid map size 30x90",
			mapSize: primitives.Size2D[uint]{Height: 30, Width: 90},
			wantErr: false,
		},
		{
			name:      "Map too small",
			mapSize:   primitives.Size2D[uint]{Height: 6, Width: 6},
			wantErr:   true,
			errorText: "game map size is too small",
		},
		{
			name:      "Map too big",
			mapSize:   primitives.Size2D[uint]{Height: 300, Width: 300},
			wantErr:   true,
			errorText: "game map size is too big",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := rand.New(rand.NewSource(randomSeedTest))
			level := NewLevelWithDefaults(source, tt.mapSize)

			err := level.Generate()

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
				// Проверяем, что всё сгенерировано
				if len(level.Rooms) == 0 {
					t.Error("Expected rooms to be generated")
				}
				if len(level.Items) == 0 {
					t.Error("Expected elixirs to be generated")
				}
				if level.Player == nil {
					t.Error("Expected player to be generated")
				}
			}
		})
	}
}

func TestLevel_GenerateRooms_Deterministic(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}

	level := NewLevelWithDefaults(source, mapSize)
	err := level.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	wantRooms := []Room{
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 8, Y: 3},
				Size:  primitives.Size2D[uint]{Width: 10, Height: 3},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 44, Y: 1},
				Size:  primitives.Size2D[uint]{Width: 16, Height: 8},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 61, Y: 3},
				Size:  primitives.Size2D[uint]{Width: 26, Height: 6},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 1, Y: 11},
				Size:  primitives.Size2D[uint]{Width: 9, Height: 8},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 40, Y: 13},
				Size:  primitives.Size2D[uint]{Width: 20, Height: 3},
			},
		},
		{
			Type: RoomTypeFinish,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 63, Y: 14},
				Size:  primitives.Size2D[uint]{Width: 5, Height: 3},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 10, Y: 24},
				Size:  primitives.Size2D[uint]{Width: 7, Height: 3},
			},
		},
		{
			Type: RoomTypeStart,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 33, Y: 23},
				Size:  primitives.Size2D[uint]{Width: 17, Height: 7},
			},
		},
		{
			Type: RoomTypeOrdinary,
			Shape: primitives.Box{
				Point: primitives.Point2D[int]{X: 66, Y: 25},
				Size:  primitives.Size2D[uint]{Width: 11, Height: 5},
			},
		},
	}
	wantFinishPortal := primitives.Box{
		Point: primitives.Point2D[int]{X: 64, Y: 15},
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
	tests := []struct {
		name         string
		wantPassages []Passage
	}{
		{
			name: "Generate Passages with deterministic rooms",
			wantPassages: []Passage{
				{
					DoorOne: primitives.Point2D[int]{X: 17, Y: 4},
					DoorTwo: primitives.Point2D[int]{X: 44, Y: 2},
					Way: []primitives.Point2D[int]{
						{X: 18, Y: 4}, {X: 19, Y: 4}, {X: 20, Y: 4}, {X: 21, Y: 4}, {X: 22, Y: 4},
						{X: 23, Y: 4}, {X: 24, Y: 4}, {X: 25, Y: 4}, {X: 26, Y: 4}, {X: 27, Y: 4},
						{X: 28, Y: 4}, {X: 29, Y: 4}, {X: 30, Y: 4}, {X: 30, Y: 3}, {X: 30, Y: 2},
						{X: 31, Y: 2}, {X: 32, Y: 2}, {X: 33, Y: 2}, {X: 34, Y: 2}, {X: 35, Y: 2},
						{X: 36, Y: 2}, {X: 37, Y: 2}, {X: 38, Y: 2}, {X: 39, Y: 2}, {X: 40, Y: 2},
						{X: 41, Y: 2}, {X: 42, Y: 2}, {X: 43, Y: 2},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 13, Y: 5},
					DoorTwo: primitives.Point2D[int]{X: 2, Y: 11},
					Way: []primitives.Point2D[int]{
						{X: 13, Y: 6}, {X: 12, Y: 6}, {X: 11, Y: 6}, {X: 10, Y: 6}, {X: 9, Y: 6},
						{X: 8, Y: 6}, {X: 7, Y: 6}, {X: 6, Y: 6}, {X: 5, Y: 6}, {X: 4, Y: 6},
						{X: 3, Y: 6}, {X: 2, Y: 6}, {X: 2, Y: 7}, {X: 2, Y: 8}, {X: 2, Y: 9},
						{X: 2, Y: 10},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 9, Y: 14},
					DoorTwo: primitives.Point2D[int]{X: 40, Y: 14},
					Way: []primitives.Point2D[int]{
						{X: 10, Y: 14}, {X: 11, Y: 14}, {X: 12, Y: 14}, {X: 13, Y: 14}, {X: 14, Y: 14},
						{X: 15, Y: 14}, {X: 16, Y: 14}, {X: 17, Y: 14}, {X: 18, Y: 14}, {X: 19, Y: 14},
						{X: 20, Y: 14}, {X: 21, Y: 14}, {X: 22, Y: 14}, {X: 23, Y: 14}, {X: 24, Y: 14},
						{X: 25, Y: 14}, {X: 26, Y: 14}, {X: 27, Y: 14}, {X: 28, Y: 14}, {X: 29, Y: 14},
						{X: 30, Y: 14}, {X: 31, Y: 14}, {X: 32, Y: 14}, {X: 33, Y: 14}, {X: 34, Y: 14},
						{X: 35, Y: 14}, {X: 36, Y: 14}, {X: 37, Y: 14}, {X: 38, Y: 14}, {X: 39, Y: 14},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 4, Y: 18},
					DoorTwo: primitives.Point2D[int]{X: 11, Y: 24},
					Way: []primitives.Point2D[int]{
						{X: 4, Y: 19}, {X: 5, Y: 19}, {X: 6, Y: 19}, {X: 7, Y: 19}, {X: 8, Y: 19},
						{X: 9, Y: 19}, {X: 10, Y: 19}, {X: 11, Y: 19}, {X: 11, Y: 20}, {X: 11, Y: 21},
						{X: 11, Y: 22}, {X: 11, Y: 23},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 16, Y: 25},
					DoorTwo: primitives.Point2D[int]{X: 33, Y: 25},
					Way: []primitives.Point2D[int]{
						{X: 17, Y: 25}, {X: 18, Y: 25}, {X: 19, Y: 25}, {X: 20, Y: 25}, {X: 21, Y: 25},
						{X: 22, Y: 25}, {X: 23, Y: 25}, {X: 24, Y: 25}, {X: 25, Y: 25}, {X: 26, Y: 25},
						{X: 27, Y: 25}, {X: 28, Y: 25}, {X: 29, Y: 25}, {X: 30, Y: 25}, {X: 31, Y: 25},
						{X: 32, Y: 25},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 49, Y: 24},
					DoorTwo: primitives.Point2D[int]{X: 66, Y: 26},
					Way: []primitives.Point2D[int]{
						{X: 50, Y: 24}, {X: 50, Y: 25}, {X: 50, Y: 26}, {X: 51, Y: 26}, {X: 52, Y: 26},
						{X: 53, Y: 26}, {X: 54, Y: 26}, {X: 55, Y: 26}, {X: 56, Y: 26}, {X: 57, Y: 26},
						{X: 58, Y: 26}, {X: 59, Y: 26}, {X: 60, Y: 26}, {X: 61, Y: 26}, {X: 62, Y: 26},
						{X: 63, Y: 26}, {X: 64, Y: 26}, {X: 65, Y: 26},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 64, Y: 16},
					DoorTwo: primitives.Point2D[int]{X: 72, Y: 25},
					Way: []primitives.Point2D[int]{
						{X: 64, Y: 17}, {X: 65, Y: 17}, {X: 66, Y: 17}, {X: 67, Y: 17}, {X: 68, Y: 17},
						{X: 69, Y: 17}, {X: 70, Y: 17}, {X: 71, Y: 17}, {X: 72, Y: 17}, {X: 72, Y: 18},
						{X: 72, Y: 19}, {X: 72, Y: 20}, {X: 72, Y: 21}, {X: 72, Y: 22}, {X: 72, Y: 23},
						{X: 72, Y: 24},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 68, Y: 8},
					DoorTwo: primitives.Point2D[int]{X: 64, Y: 14},
					Way: []primitives.Point2D[int]{
						{X: 68, Y: 9}, {X: 68, Y: 10}, {X: 68, Y: 11}, {X: 68, Y: 12},
						{X: 67, Y: 12}, {X: 66, Y: 12}, {X: 65, Y: 12}, {X: 64, Y: 12}, {X: 64, Y: 13},
					},
				},
				{
					DoorOne: primitives.Point2D[int]{X: 47, Y: 8},
					DoorTwo: primitives.Point2D[int]{X: 56, Y: 13},
					Way: []primitives.Point2D[int]{
						{X: 47, Y: 9}, {X: 47, Y: 10}, {X: 47, Y: 11}, {X: 48, Y: 11}, {X: 49, Y: 11},
						{X: 50, Y: 11}, {X: 51, Y: 11}, {X: 52, Y: 11}, {X: 53, Y: 11}, {X: 54, Y: 11},
						{X: 55, Y: 11}, {X: 56, Y: 11}, {X: 56, Y: 12},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := rand.New(rand.NewSource(randomSeedTest))
			level := NewLevelWithDefaults(source, primitives.Size2D[uint]{Height: 30, Width: 90})
			err := level.Generate()
			if err != nil {
				t.Fatalf("Generate returned an error: %v", err)
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
					t.Errorf("Passage at index %d field 'DoorTwo' mismatches: expected %v, got %v",
						i, wantPassage.DoorTwo, gotPassage.DoorTwo)
				}
				if len(gotPassage.Way) != len(wantPassage.Way) {
					t.Errorf("Passage at index %d Way slice length mismatch: expected %d, got %d", i, len(wantPassage.Way), len(gotPassage.Way))
				} else {
					for j, wantPoint := range wantPassage.Way {
						gotPoint := gotPassage.Way[j]
						if gotPoint != wantPoint {
							t.Errorf("Passage at index %d, point in Way at index %d mismatch: expected %+v, got %+v", i, j, wantPoint, gotPoint)
						}
					}
				}
			}
		})
	}
}
